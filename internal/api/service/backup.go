package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type BackupService struct {
	app *config.App
}

func NewBackupService(app *config.App) *BackupService {
	return &BackupService{app: app}
}

func (s *BackupService) CreateBackup(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateBackupRequest],
) (*connect.Response[mantraev1.CreateBackupResponse], error) {
	if err := s.app.BM.Create(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.CreateBackupResponse{}), nil
}

func (s *BackupService) ListBackups(
	ctx context.Context,
	req *connect.Request[mantraev1.ListBackupsRequest],
) (*connect.Response[mantraev1.ListBackupsResponse], error) {
	files, err := s.app.BM.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var backups []*mantraev1.Backup
	for _, file := range files {
		backups = append(backups, &mantraev1.Backup{
			Name:      file.Name,
			Size:      file.Size,
			CreatedAt: SafeTimestamp(&file.Timestamp),
		})
	}
	return connect.NewResponse(&mantraev1.ListBackupsResponse{
		Backups: backups,
	}), nil
}

func (s *BackupService) DeleteBackup(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteBackupRequest],
) (*connect.Response[mantraev1.DeleteBackupResponse], error) {
	if err := s.app.BM.Delete(ctx, req.Msg.Name); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteBackupResponse{}), nil
}

func (s *BackupService) RestoreBackup(
	ctx context.Context,
	req *connect.Request[mantraev1.RestoreBackupRequest],
) (*connect.Response[mantraev1.RestoreBackupResponse], error) {
	if !s.app.BM.IsValidBackupFile(req.Msg.Name) {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid backup file name"),
		)
	}
	if err := s.app.BM.Restore(ctx, req.Msg.Name); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.RestoreBackupResponse{}), nil
}

func (s *BackupService) DownloadBackup(
	ctx context.Context,
	req *connect.Request[mantraev1.DownloadBackupRequest],
	stream *connect.ServerStream[mantraev1.DownloadBackupResponse],
) error {
	filename := req.Msg.Name
	if req.Msg.Name == "" {
		files, err := s.app.BM.Storage.List(ctx)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
		filename = files[0].Name // Use latest backup
	}

	reader, err := s.app.BM.Storage.Retrieve(ctx, filename)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	defer reader.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return connect.NewError(connect.CodeInternal, err)
		}
		if err := stream.Send(&mantraev1.DownloadBackupResponse{
			Data: buf[:n],
		}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}

func (s *BackupService) UploadBackup(
	ctx context.Context,
	stream *connect.ClientStream[mantraev1.UploadBackupRequest],
) (*connect.Response[mantraev1.UploadBackupResponse], error) {
	// Create a pipe to stream data from gRPC into storage
	pr, pw := io.Pipe()
	defer pr.Close()

	// Start storage write in background
	filename := fmt.Sprintf("upload_%s.db", time.Now().UTC().Format("20060102_150405"))
	storeErrChan := make(chan error, 1)
	go func() {
		defer pw.Close()
		err := s.app.BM.Storage.Store(ctx, filename, pr)
		storeErrChan <- err
	}()

	for stream.Receive() {
		chunk := stream.Msg()
		if len(chunk.Data) == 0 {
			continue
		}
		if _, err := pw.Write(chunk.Data); err != nil {
			pw.CloseWithError(err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	if err := stream.Err(); err != nil {
		pw.CloseWithError(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Wait for store to finish
	if err := <-storeErrChan; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.UploadBackupResponse{}), nil
}

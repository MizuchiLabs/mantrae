package service

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
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
			CreatedAt: convert.SafeTimestamp(file.Timestamp),
		})
	}
	return connect.NewResponse(&mantraev1.ListBackupsResponse{Backups: backups}), nil
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
		if len(files) == 0 {
			// Create a new backup if none exist
			if err = s.app.BM.Create(ctx); err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
			files, err = s.app.BM.Storage.List(ctx)
			if err != nil {
				return connect.NewError(connect.CodeInternal, err)
			}
		}
		filename = files[0].Name // Use latest backup
	}

	reader, err := s.app.BM.Storage.Retrieve(ctx, filename)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	defer func() {
		if err = reader.Close(); err != nil {
			slog.Error("failed to close backup reader", "error", err)
		}
	}()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return connect.NewError(connect.CodeInternal, err)
		}
		if err := stream.Send(&mantraev1.DownloadBackupResponse{Data: buf[:n]}); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}

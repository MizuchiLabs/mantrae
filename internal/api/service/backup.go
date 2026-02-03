package service

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

type BackupService struct {
	app *config.App
}

func NewBackupService(app *config.App) *BackupService {
	return &BackupService{app: app}
}

func (s *BackupService) CreateBackup(
	ctx context.Context,
	req *mantraev1.CreateBackupRequest,
) (*mantraev1.CreateBackupResponse, error) {
	if err := s.app.BM.Create(ctx); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.CreateBackupResponse{}, nil
}

func (s *BackupService) ListBackups(
	ctx context.Context,
	req *mantraev1.ListBackupsRequest,
) (*mantraev1.ListBackupsResponse, error) {
	files, err := s.app.BM.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var backups []*mantraev1.Backup
	for _, file := range files {
		backups = append(backups, &mantraev1.Backup{
			Name:      file.Name,
			Size:      file.Size,
			CreatedAt: db.SafeTimestamp(file.Timestamp),
		})
	}
	return &mantraev1.ListBackupsResponse{Backups: backups}, nil
}

func (s *BackupService) DeleteBackup(
	ctx context.Context,
	req *mantraev1.DeleteBackupRequest,
) (*mantraev1.DeleteBackupResponse, error) {
	if err := s.app.BM.Delete(ctx, req.Name); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.DeleteBackupResponse{}, nil
}

func (s *BackupService) RestoreBackup(
	ctx context.Context,
	req *mantraev1.RestoreBackupRequest,
) (*mantraev1.RestoreBackupResponse, error) {
	if !s.app.BM.IsValidBackupFile(req.Name) {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid backup file name"),
		)
	}
	extension := strings.ToLower(filepath.Ext(req.Name))
	switch extension {
	case ".db":
		if err := s.app.BM.Restore(ctx, req.Name); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	case ".yaml", ".yml", ".json":
		s.app.BM.RestoreViaConfig(ctx, req.ProfileId, req.Name)
	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid backup file type"),
		)
	}
	return &mantraev1.RestoreBackupResponse{}, nil
}

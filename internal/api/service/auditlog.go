package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

type AuditLogService struct {
	app *config.App
}

func NewAuditLogService(app *config.App) *AuditLogService {
	return &AuditLogService{app: app}
}

func (s *AuditLogService) ListAuditLogs(
	ctx context.Context,
	req *connect.Request[mantraev1.ListAuditLogsRequest],
) (*connect.Response[mantraev1.ListAuditLogsResponse], error) {
	params := &db.ListAuditLogsParams{
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	}

	result, err := s.app.Conn.Query.ListAuditLogs(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.Query.CountAuditLogs(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	auditLogs := make([]*mantraev1.AuditLog, 0, len(result))
	for _, l := range result {
		auditLogs = append(auditLogs, l.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListAuditLogsResponse{
		AuditLogs:  auditLogs,
		TotalCount: totalCount,
	}), nil
}

package service

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
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
	var params db.ListAuditLogsParams
	if req.Msg.Limit == nil {
		params.Limit = 100
	} else {
		params.Limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *req.Msg.Offset
	}

	result, err := s.app.Conn.GetQuery().ListAuditLogs(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountAuditLogs(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.ListAuditLogsResponse{
		AuditLogs:  convert.AuditLogsToProto(result),
		TotalCount: totalCount,
	}), nil
}

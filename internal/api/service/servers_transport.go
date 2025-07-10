package service

import (
	"context"
	"errors"
	"fmt"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type ServersTransportService struct {
	app *config.App
}

func NewServersTransportService(app *config.App) *ServersTransportService {
	return &ServersTransportService{app: app}
}

func (s *ServersTransportService) GetServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.GetServersTransportRequest],
) (*connect.Response[mantraev1.GetServersTransportResponse], error) {
	var serversTransport *mantraev1.ServersTransport

	switch req.Msg.Type {
	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP:
		result, err := s.app.Conn.GetQuery().GetHttpServersTransport(ctx, req.Msg.Id)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.HTTPServersTransportToProto(&result)

	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP:
		result, err := s.app.Conn.GetQuery().GetTcpServersTransport(ctx, req.Msg.Id)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.TCPServersTransportToProto(&result)

	default:
		return nil, fmt.Errorf("invalid servers transport type")
	}

	return connect.NewResponse(
		&mantraev1.GetServersTransportResponse{ServersTransport: serversTransport},
	), nil
}

func (s *ServersTransportService) CreateServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateServersTransportRequest],
) (*connect.Response[mantraev1.CreateServersTransportResponse], error) {
	var serversTransport *mantraev1.ServersTransport
	var err error

	switch req.Msg.Type {
	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP:
		var params db.CreateHttpServersTransportParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = convert.UnmarshalStruct[schema.ServersTransport](req.Msg.Config)
		if err != nil {
			return nil, err
		}

		result, err := s.app.Conn.GetQuery().CreateHttpServersTransport(ctx, params)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.HTTPServersTransportToProto(&result)

	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP:
		var params db.CreateTcpServersTransportParams
		params.ProfileID = req.Msg.ProfileId
		params.Name = req.Msg.Name
		if req.Msg.AgentId != "" {
			params.AgentID = &req.Msg.AgentId
		}

		params.Config, err = convert.UnmarshalStruct[schema.TCPServersTransport](req.Msg.Config)
		if err != nil {
			return nil, err
		}

		result, err := s.app.Conn.GetQuery().CreateTcpServersTransport(ctx, params)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.TCPServersTransportToProto(&result)

	default:
		return nil, fmt.Errorf("invalid servers transport type")
	}

	return connect.NewResponse(
		&mantraev1.CreateServersTransportResponse{ServersTransport: serversTransport},
	), nil
}

func (s *ServersTransportService) UpdateServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateServersTransportRequest],
) (*connect.Response[mantraev1.UpdateServersTransportResponse], error) {
	var serversTransport *mantraev1.ServersTransport
	var err error

	switch req.Msg.Type {
	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP:
		var params db.UpdateHttpServersTransportParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Enabled = req.Msg.Enabled
		params.Config, err = convert.UnmarshalStruct[schema.ServersTransport](req.Msg.Config)
		if err != nil {
			return nil, err
		}

		result, err := s.app.Conn.GetQuery().UpdateHttpServersTransport(ctx, params)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.HTTPServersTransportToProto(&result)

	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP:
		var params db.UpdateTcpServersTransportParams
		params.ID = req.Msg.Id
		params.Name = req.Msg.Name
		params.Enabled = req.Msg.Enabled
		params.Config, err = convert.UnmarshalStruct[schema.TCPServersTransport](req.Msg.Config)
		if err != nil {
			return nil, err
		}

		result, err := s.app.Conn.GetQuery().UpdateTcpServersTransport(ctx, params)
		if err != nil {
			return nil, err
		}
		serversTransport = convert.TCPServersTransportToProto(&result)

	default:
		return nil, fmt.Errorf("invalid servers transport type")
	}

	return connect.NewResponse(
		&mantraev1.UpdateServersTransportResponse{ServersTransport: serversTransport},
	), nil
}

func (s *ServersTransportService) DeleteServersTransport(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteServersTransportRequest],
) (*connect.Response[mantraev1.DeleteServersTransportResponse], error) {
	switch req.Msg.Type {
	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP:
		if err := s.app.Conn.GetQuery().DeleteHttpServersTransport(ctx, req.Msg.Id); err != nil {
			return nil, err
		}

	case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP:
		if err := s.app.Conn.GetQuery().DeleteTcpServersTransport(ctx, req.Msg.Id); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid servers transport type")
	}

	return connect.NewResponse(&mantraev1.DeleteServersTransportResponse{}), nil
}

func (s *ServersTransportService) ListServersTransports(
	ctx context.Context,
	req *connect.Request[mantraev1.ListServersTransportsRequest],
) (*connect.Response[mantraev1.ListServersTransportsResponse], error) {
	var limit int64
	var offset int64
	if req.Msg.Limit == nil {
		limit = 100
	} else {
		limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		offset = 0
	} else {
		offset = *req.Msg.Offset
	}

	var serversTransports []*mantraev1.ServersTransport
	var totalCount int64

	if req.Msg.Type == nil {
		if req.Msg.AgentId == nil {
			result, err := s.app.Conn.GetQuery().
				ListServersTransportsByProfile(ctx, db.ListServersTransportsByProfileParams{
					Limit:       limit,
					Offset:      offset,
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, err
			}
			totalCount, err = s.app.Conn.GetQuery().
				CountServersTransportsByProfile(ctx, db.CountServersTransportsByProfileParams{
					ProfileID:   req.Msg.ProfileId,
					ProfileID_2: req.Msg.ProfileId,
				})
			if err != nil {
				return nil, err
			}
			serversTransports = convert.ServersTransportsByProfileToProto(
				result,
				s.app.Conn.GetQuery(),
			)
		} else {
			result, err := s.app.Conn.GetQuery().ListServersTransportsByAgent(ctx, db.ListServersTransportsByAgentParams{
				Limit:     limit,
				Offset:    offset,
				AgentID:   req.Msg.AgentId,
				AgentID_2: req.Msg.AgentId,
			})
			if err != nil {
				return nil, err
			}
			totalCount, err = s.app.Conn.GetQuery().CountServersTransportsByAgent(ctx, db.CountServersTransportsByAgentParams{
				AgentID:   req.Msg.AgentId,
				AgentID_2: req.Msg.AgentId,
			})
			if err != nil {
				return nil, err
			}
			serversTransports = convert.ServersTransportsByAgentToProto(result, s.app.Conn.GetQuery())
		}
	} else {
		var err error
		switch *req.Msg.Type {
		case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_HTTP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountHttpServersTransportsByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, err
				}
				result, err := s.app.Conn.GetQuery().ListHttpServersTransports(ctx, db.ListHttpServersTransportsParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, err
				}
				serversTransports = convert.HTTPServersTransportsToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountHttpServersTransportsByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, err
				}
				result, err := s.app.Conn.GetQuery().ListHttpServersTransportsByAgent(ctx, db.ListHttpServersTransportsByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				serversTransports = convert.HTTPServersTransportsToProto(result)
			}

		case mantraev1.ServersTransportType_SERVERS_TRANSPORT_TYPE_TCP:
			if req.Msg.AgentId == nil {
				totalCount, err = s.app.Conn.GetQuery().CountTcpServersTransportsByProfile(ctx, req.Msg.ProfileId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpServersTransports(ctx, db.ListTcpServersTransportsParams{
					ProfileID: req.Msg.ProfileId,
					Limit:     limit,
					Offset:    offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				serversTransports = convert.TCPServersTransportsToProto(result)
			} else {
				totalCount, err = s.app.Conn.GetQuery().CountTcpServersTransportsByAgent(ctx, req.Msg.AgentId)
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				result, err := s.app.Conn.GetQuery().ListTcpServersTransportsByAgent(ctx, db.ListTcpServersTransportsByAgentParams{
					AgentID: req.Msg.AgentId,
					Limit:   limit,
					Offset:  offset,
				})
				if err != nil {
					return nil, connect.NewError(connect.CodeInternal, err)
				}
				serversTransports = convert.TCPServersTransportsToProto(result)
			}

		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid servers transport type"))
		}

		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.ListServersTransportsResponse{
		ServersTransports: serversTransports,
		TotalCount:        totalCount,
	}), nil
}

package service

import (
	"context"

	"connectrpc.com/connect"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
)

type TraefikInstanceService struct {
	app *config.App
}

func NewTraefikInstanceService(app *config.App) *TraefikInstanceService {
	return &TraefikInstanceService{app: app}
}

func (s *TraefikInstanceService) GetTraefikInstance(
	ctx context.Context,
	req *connect.Request[mantraev1.GetTraefikInstanceRequest],
) (*connect.Response[mantraev1.GetTraefikInstanceResponse], error) {
	result, err := s.app.Conn.Query.GetTraefikInstanceByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&mantraev1.GetTraefikInstanceResponse{
		TraefikInstance: result.ToProto(),
	}), nil
}

func (s *TraefikInstanceService) DeleteTraefikInstance(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteTraefikInstanceRequest],
) (*connect.Response[mantraev1.DeleteTraefikInstanceResponse], error) {
	instance, err := s.app.Conn.Query.GetTraefikInstanceByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}
	if err := s.app.Conn.Query.DeleteTraefikInstance(ctx, req.Msg.Id); err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_TraefikInstance{
			TraefikInstance: instance.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.DeleteTraefikInstanceResponse{}), nil
}

func (s *TraefikInstanceService) ListTraefikInstances(
	ctx context.Context,
	req *connect.Request[mantraev1.ListTraefikInstancesRequest],
) (*connect.Response[mantraev1.ListTraefikInstancesResponse], error) {
	params := &db.ListTraefikInstancesParams{
		ProfileID: req.Msg.ProfileId,
		Limit:     req.Msg.Limit,
		Offset:    req.Msg.Offset,
	}

	result, err := s.app.Conn.Query.ListTraefikInstances(ctx, params)
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Query.CountTraefikInstances(ctx, req.Msg.ProfileId)
	if err != nil {
		return nil, err
	}

	instances := make([]*mantraev1.TraefikInstance, 0, len(result))
	for _, i := range result {
		instances = append(instances, i.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListTraefikInstancesResponse{
		TraefikInstances: instances,
		TotalCount:       totalCount,
	}), nil
}

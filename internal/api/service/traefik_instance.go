package service

import (
	"context"

	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/store/db"
)

type TraefikInstanceService struct {
	app *config.App
}

func NewTraefikInstanceService(app *config.App) *TraefikInstanceService {
	return &TraefikInstanceService{app: app}
}

func (s *TraefikInstanceService) GetTraefikInstance(
	ctx context.Context,
	req *mantraev1.GetTraefikInstanceRequest,
) (*mantraev1.GetTraefikInstanceResponse, error) {
	result, err := s.app.Conn.Q.GetTraefikInstanceByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mantraev1.GetTraefikInstanceResponse{TraefikInstance: result.ToProto()}, nil
}

func (s *TraefikInstanceService) DeleteTraefikInstance(
	ctx context.Context,
	req *mantraev1.DeleteTraefikInstanceRequest,
) (*mantraev1.DeleteTraefikInstanceResponse, error) {
	instance, err := s.app.Conn.Q.GetTraefikInstanceByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if err := s.app.Conn.Q.DeleteTraefikInstance(ctx, req.Id); err != nil {
		return nil, err
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_TraefikInstance{
			TraefikInstance: instance.ToProto(),
		},
	})
	return &mantraev1.DeleteTraefikInstanceResponse{}, nil
}

func (s *TraefikInstanceService) ListTraefikInstances(
	ctx context.Context,
	req *mantraev1.ListTraefikInstancesRequest,
) (*mantraev1.ListTraefikInstancesResponse, error) {
	params := &db.ListTraefikInstancesParams{
		ProfileID: req.ProfileId,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	result, err := s.app.Conn.Q.ListTraefikInstances(ctx, params)
	if err != nil {
		return nil, err
	}
	totalCount, err := s.app.Conn.Q.CountTraefikInstances(ctx, req.ProfileId)
	if err != nil {
		return nil, err
	}

	instances := make([]*mantraev1.TraefikInstance, 0, len(result))
	for _, i := range result {
		instances = append(instances, i.ToProto())
	}
	return &mantraev1.ListTraefikInstancesResponse{
		TraefikInstances: instances,
		TotalCount:       totalCount,
	}, nil
}

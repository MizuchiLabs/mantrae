package service

import (
	"context"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/events"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type EventService struct {
	app *config.App
}

func NewEventService(app *config.App) *EventService {
	return &EventService{app: app}
}

func (s *EventService) StreamProfileEvents(
	ctx context.Context,
	req *connect.Request[mantraev1.StreamProfileEventsRequest],
	stream *connect.ServerStream[mantraev1.ProfileEvent],
) error {
	// Create filtered event channel
	eventChan := make(chan *mantraev1.ProfileEvent, 100)

	// Register with broadcaster
	filter := &events.EventFilter{
		ProfileID:     req.Msg.ProfileId,
		ResourceTypes: req.Msg.ResourceTypes,
	}

	s.app.Event.RegisterProfileClient(filter, eventChan)
	defer s.app.Event.UnregisterProfileClient(eventChan)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-eventChan:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (s *EventService) StreamGlobalEvents(
	ctx context.Context,
	req *connect.Request[mantraev1.StreamGlobalEventsRequest],
	stream *connect.ServerStream[mantraev1.GlobalEvent],
) error {
	// Verify admin permissions
	// user := s.app.GetUserFromContext(ctx)
	// if !user.IsAdmin {
	// 	return connect.NewError(connect.CodePermissionDenied, nil)
	// }

	eventChan := make(chan *mantraev1.GlobalEvent, 100)
	filter := &events.GlobalEventFilter{ResourceTypes: req.Msg.ResourceTypes}

	s.app.Event.RegisterGlobalClient(filter, eventChan)
	defer s.app.Event.UnregisterGlobalClient(eventChan)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-eventChan:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

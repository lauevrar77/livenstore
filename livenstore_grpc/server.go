package livenstore_grpc

import (
	context "context"
	"time"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/services"
)

type Server struct {
	UnimplementedLivenstoreServer

	ES services.EventStore
}

func (s *Server) Publish(c context.Context, req *PublishEventRequest) (*PublishEventReply, error) {
	e := domain.Event{
		ID:        ulid.Make(),
		Type:      req.Type,
		Timestamp: uint64(time.Now().Unix()),
		Data:      req.Data,
	}
	_, err := s.ES.PublishEvent(e)
	if err != nil {
		return nil, err
	}
	return &PublishEventReply{Id: e.ID.String()}, nil
}

func (s *Server) EventByID(c context.Context, req *EventByIDRequest) (*EventResponse, error) {
	ulid, err := ulid.Parse(req.Id)
	if err != nil {
		return nil, err
	}
	e, err := s.ES.EventByID(ulid)
	if err != nil {
		return nil, err
	}
	return &EventResponse{
		Event: &Event{
			Id:        e.ID.String(),
			Type:      e.Type,
			Timestamp: e.Timestamp,
			Data:      e.Data,
		},
	}, nil
}

func (s *Server) LinkEventToStream(c context.Context, req *LinkEventToStreamRequest) (*EmptyResponse, error) {
	ulid, err := ulid.Parse(req.EventId)
	if err != nil {
		return nil, err
	}
	err = s.ES.LinkToStream(req.StreamName, ulid)
	if err != nil {
		return nil, err
	}
	return &EmptyResponse{}, nil
}

func (s *Server) ReadStream(c context.Context, req *ReadStreamRequest) (*Stream, error) {
	stream, err := s.ES.ReadStream(req.StreamName)
	if err != nil {
		return nil, err
	}
	eventIDs := make([]string, len(stream.EventIDs))
	for i, id := range stream.EventIDs {
		eventIDs[i] = id.String()
	}
	return &Stream{
		Name:     stream.Name,
		EventIds: eventIDs,
	}, nil
}

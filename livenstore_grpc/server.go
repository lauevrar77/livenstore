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
	_, err := s.ES.WriteEvent(e)
	if err != nil {
		return nil, err
	}
	return &PublishEventReply{Id: e.ID.String()}, nil
}

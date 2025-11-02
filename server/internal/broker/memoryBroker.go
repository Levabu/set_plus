package broker

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type MemoryBroker struct {
	onReceiveEventCallback func(roomID uuid.UUID, event domain.Event) error
}

func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{
	}
}

func (p *MemoryBroker) SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, msg any)) error {
	return nil
}

func (s *MemoryBroker) PublishRoomUpdate(ctx context.Context, roomID uuid.UUID, event domain.Event) error {
	if s.onReceiveEventCallback == nil {
		return nil
	}
	// only emulates publishing, so handle immediately
	go func() {
		s.onReceiveEventCallback(roomID, event)
	}()
	return nil
}

func (s *MemoryBroker) SetEventCallback(callback func(roomID uuid.UUID, event domain.Event) error) {
	s.onReceiveEventCallback = callback
} 

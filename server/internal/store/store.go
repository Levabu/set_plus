package store

import (
	"context"
	"server/internal/domain"
	"server/internal/game"

	"github.com/google/uuid"
)

type Store interface {
	SetRoom(ctx context.Context, room *domain.Room) error
	GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error)

	SetGameState(ctx context.Context, game *game.Game) error
	GetGameState(ctx context.Context, id uuid.UUID) (*game.Game, error)

	PublishRoomUpdate(ctx context.Context, roomID uuid.UUID, event domain.Event) error
}
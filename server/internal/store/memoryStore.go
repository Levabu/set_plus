package store

import (
	"context"
	"errors"
	"server/internal/domain"
	"server/internal/game"
	"sync"

	"github.com/google/uuid"
)

type MemoryStore struct {
	games map[uuid.UUID]*game.Game
	rooms map[uuid.UUID]*domain.Room
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games: make(map[uuid.UUID]*game.Game),
		rooms: make(map[uuid.UUID]*domain.Room),
	}
}

func (s *MemoryStore) SetRoom(ctx context.Context, room *domain.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rooms[room.ID] = room
	return nil
}
func (s *MemoryStore) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	room, ok := s.rooms[id]
	if !ok {
		err = errors.New("room doesn't exist")
	}
	return room, err
}

func (s *MemoryStore) SetGameState(ctx context.Context, game *game.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[game.GameID] = game
	return nil
}
func (s *MemoryStore) GetGameState(ctx context.Context, id uuid.UUID) (*game.Game, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var err error
	gameState, ok := s.games[id]
	if !ok {
		err = errors.New("game doesn't exist")
	}
	return gameState, err
}

func (s *MemoryStore) CleanupAfterGame(ctx context.Context, gameID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.games, gameID)
}

func (s *MemoryStore) CleanupStoreRoom(ctx context.Context, roomID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rooms, roomID)
}

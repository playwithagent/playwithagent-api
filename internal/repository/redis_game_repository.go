package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"playwithagent-xo/internal/game"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)


type RedisGameRepository struct {
	rdb *redis.Client
	idMutex *sync.Mutex
	nextGameID int
}

func NewRedisGameRepository(rdb *redis.Client) *RedisGameRepository {
	return &RedisGameRepository{
		rdb: rdb,
		idMutex: &sync.Mutex{},
		nextGameID: 1,
	}
}

func (r *RedisGameRepository) GetNextGameID() int {
	r.idMutex.Lock()
	defer r.idMutex.Unlock()

	nextID := r.nextGameID
	r.nextGameID++
	return nextID
}

func (r *RedisGameRepository) SaveActiveGame(ctx context.Context, g *game.Game) (int, error) {

	if g.ID == 0 {
		g.ID = r.GetNextGameID()
	}

	gameID := g.ID

	redisKey := fmt.Sprintf("game:%d", gameID)
	gameJSON, err := json.Marshal(g)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal game: %w", err)
	}

	err = r.rdb.Set(ctx, redisKey, gameJSON, time.Hour*1).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to save game to Redis: %w", err)
	}
	fmt.Printf("Stored game %d in Redis. Key: %s\n", gameID, redisKey)
	return gameID, nil
}
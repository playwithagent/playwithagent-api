package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"playwithagent-xo/internal/game"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("resource not found") 



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

func (r *RedisGameRepository) GetActiveGame(ctx context.Context, gameID int) (*game.Game, error) {
	redisKey := fmt.Sprintf("game:%d", gameID)
	gameJSON, err := r.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil{
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get game from Redis: %w", err)
	}

	var g game.Game
	err = json.Unmarshal([]byte(gameJSON), &g)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal game: %w", err)
	}
	return &g, nil
}
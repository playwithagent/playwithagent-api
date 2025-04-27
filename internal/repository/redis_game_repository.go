package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"playwithagent-xo/internal/game"
	"time"

	"github.com/go-redis/redis/v8"
)

var ErrNotFound = errors.New("resource not found") 

const redisGameIDCounterKey = "next_game_id"


type RedisGameRepository struct {
	rdb *redis.Client
}

func NewRedisGameRepository(rdb *redis.Client) *RedisGameRepository {
	return &RedisGameRepository{
		rdb: rdb,
	}
}

// func (r *RedisGameRepository) GetNextGameID() int {
// 	r.idMutex.Lock()
// 	defer r.idMutex.Unlock()

// 	nextID := r.nextGameID
// 	r.nextGameID++
// 	return nextID
// }

func (r *RedisGameRepository) SaveActiveGame(ctx context.Context, g *game.Game) (int, error) {

	genearatedGameID, err := r.rdb.Incr(ctx, redisGameIDCounterKey).Result()

	if err != nil {
		return 0, fmt.Errorf("failed to generate next game ID from Redis: %w", err)
	}
	gameID := int(genearatedGameID)
	g.ID = gameID

	redisKey := fmt.Sprintf("game:%d", gameID)
	gameJSON, err := json.Marshal(g)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal game with generated ID %d: %w", gameID, err)
	}

	err = r.rdb.Set(ctx, redisKey, gameJSON, time.Hour*1).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to save game %d to Redis: %w", gameID, err)
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
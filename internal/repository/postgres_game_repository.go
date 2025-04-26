package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"playwithagent-xo/internal/game"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresGameRepository struct {
	db *pgxpool.Pool
}

func NewPostgresGameRepository(db *pgxpool.Pool) *PostgresGameRepository {
	return &PostgresGameRepository{
		db: db,
	}
}

func (r *PostgresGameRepository) SaveCompletedGame(ctx context.Context, g *game.Game) error {
	if g == nil {
		return errors.New("can not save nil game")
	}
	if g.GameStatus != game.GameInProgress {
		return fmt.Errorf("cannot save game %d to completed table: game is still in progress", g.ID)
	}
	if g.EndTime.IsZero() {
		return fmt.Errorf("cannot save game %d to completed table: EndTime is not set", g.ID)
	}
	finalBoardStateJSON, err := json.Marshal(g.Board.Cells)	
	if err != nil {
		return fmt.Errorf("error marshalling final board state for game %d: %w", g.ID, err)
	}

	sql := `
        INSERT INTO completed_games (game_id, start_time, end_time, status, final_board_state)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (game_id) DO NOTHING; -- Optional: Prevent errors if trying to insert duplicate ID
    `

	commandTag, err := r.db.Exec(ctx, sql,
		g.ID,
		g.StartTime,
		g.EndTime,
		int(g.GameStatus),
		finalBoardStateJSON,
	)
	if err != nil {
		return fmt.Errorf("error saving game %d to completed table: %w", g.ID, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected when saving game %d to completed table", g.ID)
	} else {
		fmt.Printf("Game %d saved to completed table successfully\n", g.ID)
	}
	
	return nil
}
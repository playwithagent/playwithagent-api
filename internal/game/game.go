package game

import (
	"errors"
	"fmt"
	"time"
)


type GameStatus int
const (
	GameInProgress GameStatus = iota
	PlayerXWon
	PlayerOWon
	Draw
)

type Game struct {
    ID            int
    Board         *Board
    CurrentPlayer Player
    StartTime     time.Time
    EndTime       time.Time
    GameStatus    GameStatus
}

func NewGame(id int, startingPlayer Player) *Game {
	return &Game{
		ID: id,
		Board: NewBoard(),
		CurrentPlayer: startingPlayer,
		StartTime: time.Now(),
		GameStatus: GameInProgress,
	}
}


func (g *Game) MakeMove(row, col int) error{
	if g.GameStatus != GameInProgress {
		return fmt.Errorf("can not make move: game stauts is %v", g.GameStatus)
	}

	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		return fmt.Errorf("invalid move: position (%d, %d) is out of bounds", row, col)
	}

	if g.Board.Get(row, col) != EmptyPlayer {
		return fmt.Errorf("invalid move: cell (%d, %d) already occupied", row, col)
	}
	currentPlayer := g.CurrentPlayer
	err := g.Board.Set(row, col, currentPlayer)
	if err != nil {
		return err
	}

	if g.CheckWin() {
		if currentPlayer == PlayerX {
			g.GameStatus = PlayerXWon
			} else {
			g.GameStatus = PlayerOWon
		}
		g.EndTime = time.Now()
		return nil
	} else if g.CheckDraw() {
		g.GameStatus = Draw
		g.EndTime = time.Now()
		return nil
	}

	err = g.SwitchTurn()
	if err != nil {
		return fmt.Errorf("failed to switch turn: %w", err)
	}
	return nil
	

}

func (g *Game) CheckWin() bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if g.Board.Get(i, 0) == g.CurrentPlayer && g.Board.Get(i, 1) == g.CurrentPlayer && g.Board.Get(i, 2) == g.CurrentPlayer {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if g.Board.Get(0, i) == g.CurrentPlayer && g.Board.Get(1, i) == g.CurrentPlayer && g.Board.Get(2, i) == g.CurrentPlayer {
			return true
		}
	}

	// Check diagonals
	if (g.Board.Get(0, 0) == g.CurrentPlayer && g.Board.Get(1, 1) == g.CurrentPlayer && g.Board.Get(2, 2) == g.CurrentPlayer) ||
	   (g.Board.Get(0, 2) == g.CurrentPlayer && g.Board.Get(1, 1) == g.CurrentPlayer && g.Board.Get(2, 0) == g.CurrentPlayer) {
		return true
	}

	return false
}

func (g *Game) CheckDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board.Get(i, j) == EmptyPlayer {
				return false
			}
		}
	}
	return true
}

func (g *Game) SwitchTurn() error {
	if g.GameStatus != GameInProgress {
		return errors.New("cannot switch turn: game is not in progress")
	}
	if g.CurrentPlayer == PlayerX {
		g.CurrentPlayer = PlayerO
	} else {
		g.CurrentPlayer = PlayerX
	}
	return nil
}

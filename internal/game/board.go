package game

import "fmt"

type Board struct {
	Cells [9]Player
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board)Get(row, col int) Player {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		panic("Invalid row or column")
	}
	return b.Cells[row*3+col]
}

func (b *Board) Set(row, col int, player Player) error {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return fmt.Errorf("invalid move: position (%d, %d) out of bounds", row, col)
	}
	index := row*3 + col
	if b.Cells[index] != EmptyPlayer {
		return fmt.Errorf("invalid move: cell (%d, %d) already occupied", row, col)
	}

	if player == EmptyPlayer {
		return fmt.Errorf("invalid operation: cannot set cell (%d, %d) back to empty explicitly", row, col) 
	}

	b.Cells[index] = player
	return nil
}

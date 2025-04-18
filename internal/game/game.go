package game

type Player int

const (
	EmptyPlayer	Player = iota
	PlayerX
	PlayerO
)	

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

func (b *Board) Set(row, col int, player Player) {
	if row < 0 || row > 2 || col < 0 || col > 2 {
		panic("Invalid row or column")
	}
	
	if b.Cells[row*3 + col] != EmptyPlayer {
		panic("Cell already occupied")
	}

	b.Cells[row*3 + col] = player
}




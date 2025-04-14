package entity

import "time"

type GameStatus string
const (
	Running GameStatus = "Running"
	Finished GameStatus = "Finished"

)
type Game struct {
	ID int
	Board [][]int
	currentPlayer int
	startTime time.Time
	GameStatus GameStatus
	PlayersHuman []User
	PlayersAI []AIPlayer
	Winner int
}

type User struct {
	ID int
}

type AIPlayer struct {
	ID int
	Name string
}
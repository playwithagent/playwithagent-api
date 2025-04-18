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
	CurrentPlayer PlayerNmber
	StartTime time.Time
	GameStatus GameStatus
	PlayersHuman User
	PlayersAI AIPlayer
	Winner int
}

type User struct {
	ID int
}

type AIPlayer struct {
	ID int
	Name string
}

type PlayerNmber int 
const (
	Player1 PlayerNmber = 1
	Player2 PlayerNmber = 2
)




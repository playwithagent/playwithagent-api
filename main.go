package main

import (
	"fmt"
	"playwithagent-xo/cmd/httpserver"
	"playwithagent-xo/config"
)



func CheckWinner(board [][]int) (int ,bool) {
	for i := 0; i < 3; i++ {
		// Check rows
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] != 0 {
			return board[i][0], true
		}
		// Check columns
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] != 0 {
			return board[0][i], true
		}
	}
	// Check diagonals
	if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != 0 {
		return board[0][0], true
	}

	if board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != 0 {
		return board[0][2], true
	}
	return 0, false
}

func BoardEmpty(board [][]int) bool {
	for _, row := range board {
		for _, cell := range row {
			if cell == 0 {
				return true
			}
		}
	}
	return false		
}

func GameMove(board [][]int, player int, row int, col int) bool {
	if board[row][col] == 0 {
		board[row][col] = player
		return true
	}
	return false
}

func main(){


	cfg := config.Config{
		HTTPServer: config.HTTPServer{
			Port: 8080,},
	}

	server := httpserver.NewServer(cfg)
	fmt.Println("Start Echo server")

	server.Serve()

	// Initialize the game
	// This is a simple representation of a game board
	// game := entity.Game{
	// 	ID: 1,
	// 	Board: [][]int{
	// 		{0, 0, 0},
	// 		{0, 0, 0},
	// 		{0, 0, 0},
	// 	},
	// 	CurrentPlayer: entity.Player1,
	// 	StartTime: time.Now(),
	// 	GameStatus: entity.Running,
	// 	PlayersHuman: entity.User{ID: 1},
	// 	PlayersAI: entity.AIPlayer{ID: 2, Name: "AI"},
	// 	Winner: 0,
	// }
	// // Simulate game play
	// for BoardEmpty(game.Board){

	// 	winner, isWinner := CheckWinner(game.Board)
	// 	if isWinner {
	// 		game.Winner = winner
	// 		game.GameStatus = entity.Finished
	// 		break
	// 	}

	// 	// Simulate a move
		
	// 	isMove := GameMove(game.Board, game.CurrentPlayer, 0, 0)
	// 	if isMove {
	// 		game.CurrentPlayer = entity.Player2
	// 	} else {
	// 		game.GameStatus = entity.InvalidMove
	// 		break
	// 	}


	// }

	
}
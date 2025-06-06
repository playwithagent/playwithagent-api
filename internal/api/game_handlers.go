package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"playwithagent-xo/internal/game"
	"playwithagent-xo/internal/repository"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ActiveGameRepository interface {
	SaveActiveGame(ctx context.Context, g *game.Game) (int, error)
	GetActiveGame(ctx context.Context, gameID int) (*game.Game, error)
}

type GameRepository interface {
	SaveCompletedGame(ctx context.Context, g *game.Game) error
}
type GameHandler struct {
	activeGameRepo ActiveGameRepository
	completedGameRepo GameRepository
}

func NewGameHandler(activeRepo ActiveGameRepository, completedRepo GameRepository) *GameHandler {
	return &GameHandler{activeGameRepo: activeRepo, completedGameRepo: completedRepo}
}

func (h *GameHandler) CreateGame(c echo.Context) error {

	newGame := game.NewGame(0, game.PlayerX)

	gameID, err := h.activeGameRepo.SaveActiveGame(c.Request().Context(), newGame)
	if err != nil {
		c.Logger().Errorf("Failed to store game in repository: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create game")
	}
	newGame.ID = gameID

	locationURL := fmt.Sprintf("/api/v1/games/%d", gameID)
	c.Response().Header().Set(echo.HeaderLocation, locationURL)
	return c.JSON(http.StatusCreated, newGame)
}

func (h *GameHandler) GetGame(c echo.Context) error {
	gameID := c.Param("game_id")
	
	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid Game ID format: %s", gameID))
	
	}
	if gameIDInt <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Game ID must be a positive integer")
	}
	fmt.Printf("Converted game_id to integer: %d\n", gameIDInt)
	gameSaved, err := h.activeGameRepo.GetActiveGame(c.Request().Context(), gameIDInt)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Game with ID %d not found", gameIDInt))
		}
		c.Logger().Errorf("Failed to retrieve game %d: %v", gameIDInt, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve game information")
	}

	return c.JSON(http.StatusOK, gameSaved)	
}

func (h *GameHandler) MakeMoveHandler(c echo.Context) error {
	gameIDStr := c.Param("game_id")
	
	
	gameIDInt, err := strconv.Atoi(gameIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid Game ID format: %s", gameIDStr))
	}
	if gameIDInt <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Game ID must be a positive integer")
	}
	var req MakeMoveRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
	}

	
	rowInt, colInt := req.Row, req.Col
	if rowInt < 0 || rowInt > 2 || colInt < 0 || colInt > 2 {
		return echo.NewHTTPError(http.StatusBadRequest, "Row and column must be between 0 and 2")
	}

	gameSaved, err := h.activeGameRepo.GetActiveGame(c.Request().Context(), gameIDInt)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Game with ID %d not found", gameIDInt))
		}
		c.Logger().Errorf("Failed to retrieve game %d: %v", gameIDInt, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve game information")
	}

	if gameSaved.GameStatus != game.GameInProgress {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Game is not in progress. Current status: %v", gameSaved.GameStatus))
	}
	err = gameSaved.MakeMove(rowInt, colInt)
	if err != nil {
		c.Logger().Errorf("Failed to make move: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid move: %v", err))
	}
	_, err = h.activeGameRepo.SaveActiveGame(c.Request().Context(), gameSaved)
	if err != nil {
		c.Logger().Errorf("CRITICAL: Failed to save updated game state %d after move: %v", gameIDInt, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Move successful but failed to save updated game state.")
	}
	if gameSaved.GameStatus != game.GameInProgress {
		
		err = h.completedGameRepo.SaveCompletedGame(c.Request().Context(), gameSaved)
		if err != nil {
			c.Logger().Errorf("Failed to save completed game %d: %v", gameIDInt, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save completed game.")
		}
	}


	return c.JSON(http.StatusOK, gameSaved)
}
	
type MakeMoveRequest struct {
	Row 		 int    `json:"row"`
	Col 		 int    `json:"col"`
}
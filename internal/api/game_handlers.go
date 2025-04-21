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

type GameRepository interface {
	SaveActiveGame(ctx context.Context, g *game.Game) (int, error)
	GetActiveGame(ctx context.Context, gameID int) (*game.Game, error)
}

type GameHandler struct {
	repo GameRepository
}

func NewGameHandler(repo GameRepository) *GameHandler {
	return &GameHandler{repo: repo}
}

func (h *GameHandler) CreateGame(c echo.Context) error {

	newGame := game.NewGame(0, game.PlayerX)

	gameID, err := h.repo.SaveActiveGame(c.Request().Context(), newGame)
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
	gameSaved, err := h.repo.GetActiveGame(c.Request().Context(), gameIDInt)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Game with ID %d not found", gameIDInt))
		}
		c.Logger().Errorf("Failed to retrieve game %d: %v", gameIDInt, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve game information")
	}

	return c.JSON(http.StatusOK, gameSaved)	
}
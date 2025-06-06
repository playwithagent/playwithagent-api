package api

import (
	"context"
	"fmt"
	"net/http"
	"playwithagent-xo/internal/game"

	"github.com/labstack/echo/v4"
)

type GameRepository interface {
	SaveActiveGame(ctx context.Context, g *game.Game) (int, error)
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
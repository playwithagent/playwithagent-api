package httpserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) healthCheck(c echo.Context) error {

    return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})	
}
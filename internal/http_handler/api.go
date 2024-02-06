package http_handler

import (
	global "async_course/main"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterAPI(g *echo.Group) {
	g.GET("/hello", h.hello)
	g.POST("/add-user", h.addUser)
}

func (h *Handler) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) addUser(c echo.Context) error {
	payload, err := validatePayload[global.AddUserReq](c)
	if err != nil {
		return err
	}
	err = h.s.AddUser(c.Request().Context(), payload.UserId)
	if err != nil {
		slog.Error("error while adding new user", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error while adding user")
	}
	return c.String(http.StatusOK, "Added user")
}

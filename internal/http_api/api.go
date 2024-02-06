package http_handler

import (
	service "async_course/main/internal/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type HttpAPI struct {
	config *viper.Viper
	s      *service.Service
}

func NewHttpAPI(config *viper.Viper, service *service.Service) *HttpAPI {
	return &HttpAPI{
		config: config,
		s:      service,
	}
}

func validatePayload[T any](c echo.Context) (T, error) {
	var payload T
	if err := c.Bind(&payload); err != nil {
		return payload, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return payload, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return payload, nil
}

package handlers

import (
	"homeserver/services"

	"github.com/labstack/echo/v4"
)

type ProxySubHandler struct {
	service *services.ProxySubService
}

func NewProxySubHandler(service *services.ProxySubService) *ProxySubHandler {
	return &ProxySubHandler{
		service,
	}
}

func (s *ProxySubHandler) RegisterRouter(group *echo.Group) error {
	group.GET("/clash2surge", s.HandleClashToSurge)
	return nil
}

func (s *ProxySubHandler) GetAPIPrefix() string {
	return "/api/proxysub"
}

func (s *ProxySubHandler) GetMiddlewareFuncs() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *ProxySubHandler) HandleClashToSurge(ctx echo.Context) error {
	return nil
}

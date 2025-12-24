package services

import "github.com/labstack/echo/v4"

type ProxySubService struct{}

func NewProxySubService() *ProxySubService {
	return &ProxySubService{}
}

func (s *ProxySubService) ClashToSurge(ctx echo.Context) (string, error) {

}

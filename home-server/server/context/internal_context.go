package context

import (
	"context"
	"homeserver/common"
	"homeserver/infrastructure"

	"github.com/labstack/echo/v4"
)

type GContext struct {
	echo.Context
}

func (s *GContext) GetContext() context.Context {
	return s.Request().Context()
}

func (s *GContext) GetAppConfig() common.Config {
	return common.GetAppConfig()
}

func (s *GContext) GetGatewaySSHClient() (*infrastructure.SSHClientWrapper, error) {
	config := s.GetAppConfig()
	return infrastructure.GetSSHClient(context.Background(), config.GatewaySSHUser, config.GatewaySSHAddress, []byte(config.SSHPrivateKey))
}

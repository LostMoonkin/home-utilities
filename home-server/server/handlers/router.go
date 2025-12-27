package handlers

import (
	"homeserver/common"
	"sync"

	"github.com/labstack/echo/v4"
)

type HandlerRegister interface {
	RegisterRouter(group *echo.Group) error
	GetAPIPrefix() string
	GetMiddlewareFuncs() []echo.MiddlewareFunc
}

type HandlerRegisters []HandlerRegister

var handlerInitOnce sync.Once
var handlerRegisterHolder []HandlerRegister = []HandlerRegister{}

func RegisterHandler(h HandlerRegister) {
	handlerRegisterHolder = append(handlerRegisterHolder, h)
}

func SetupRouter(app *echo.Echo) {
	handlerInitOnce.Do(func() {
		for _, register := range handlerRegisterHolder {
			common.Log.Info().
				Str("API Prefix", register.GetAPIPrefix()).
				Int("MiddlewareFuncs nums", len(register.GetMiddlewareFuncs())).
				Msg("Register router handlers")
			err := register.RegisterRouter(app.Group(register.GetAPIPrefix(), register.GetMiddlewareFuncs()...))
			if err != nil {
				common.Log.Fatal().Err(err).Msg("Register router handler fatal error.")
				panic("Register router handler fatal error.")
			}
		}
	})
}

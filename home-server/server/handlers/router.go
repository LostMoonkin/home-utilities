package handlers

import (
	"homeserver/common"
	"homeserver/context"
	"sync"

	"github.com/labstack/echo/v4"
)

type HandlerRegister interface {
	RegisterRouter(group *echo.Group) error
	GetAPIPrefix() string
	GetMiddlewareFunc() []echo.MiddlewareFunc
}

type ContextHandlerFunc func(ctx context.GContext) error
type HandlerRegisters []HandlerRegister

var handlerInitOnce sync.Once
var handlerRegisterHolder []HandlerRegister

func RegisterHandler(h HandlerRegister) {
	handlerRegisterHolder = append(handlerRegisterHolder, h)
}

func SetupRouter(app *echo.Echo) {
	handlerInitOnce.Do(func() {
		for _, register := range handlerRegisterHolder {
			common.Log.Info().
				Str("API Prefix", register.GetAPIPrefix()).
				Int("MiddlewareFunc nums", len(register.GetMiddlewareFunc())).
				Msg("Register router handlers")
			err := register.RegisterRouter(app.Group(register.GetAPIPrefix(), register.GetMiddlewareFunc()...))
			if err != nil {
				common.Log.Fatal().Err(err).Msg("Register router handler fatal error.")
				panic("Register router handler fatal error.")
			}
		}
	})
}

func wrapHandlerContext(handlerFunc ContextHandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		gCtx := ctx.(context.GContext)
		return handlerFunc(gCtx)
	}
}

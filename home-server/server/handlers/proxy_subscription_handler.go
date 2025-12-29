package handlers

import (
	"homeserver/common"
	"homeserver/context"
	"homeserver/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Clash2SurgeParam struct {
	typeFilters []string
}

type ProxySubHandler struct {
	service *services.ProxySubService
}

func NewProxySubHandler(service *services.ProxySubService) *ProxySubHandler {
	return &ProxySubHandler{
		service,
	}
}

func (s *ProxySubHandler) RegisterRouter(group *echo.Group) error {
	group.GET("/clash2surge", wrapHandlerContext(s.HandleClashToSurge))
	return nil
}

func (s *ProxySubHandler) GetAPIPrefix() string {
	return "/api/proxysub"
}

func (s *ProxySubHandler) GetMiddlewareFunc() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *ProxySubHandler) HandleClashToSurge(ctx context.GContext) error {
	param := &Clash2SurgeParam{}
	err := echo.QueryParamsBinder(ctx).Strings("type", &param.typeFilters).BindError()
	if err != nil {
		common.Log.Warn().Err(err).Str("param", ctx.Request().URL.RawQuery).Msg("Bind query param error.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	surgeSubText, err := s.service.ClashToSurge(ctx, param.typeFilters)
	if err != nil {
		common.Log.Error().Err(err).Msg("HandleClashToSurge error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.String(http.StatusOK, surgeSubText)
}

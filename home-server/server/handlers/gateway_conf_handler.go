package handlers

import (
	"homeserver/common"
	"homeserver/context"
	"homeserver/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetContentParam struct {
	names []string
}
type GatewayConfHandler struct {
	service *services.GatewayConfService
}

func NewGatewayConfHandler(service *services.GatewayConfService) *GatewayConfHandler {
	return &GatewayConfHandler{service: service}
}

func (s *GatewayConfHandler) RegisterRouter(group *echo.Group) error {
	group.GET("/listconf", wrapHandlerContext(s.ListAllConfigs))
	group.GET("/confdetail", wrapHandlerContext(s.GetConfContent))
	return nil
}

func (s *GatewayConfHandler) GetAPIPrefix() string {
	return "/api/gateway"
}

func (s *GatewayConfHandler) GetMiddlewareFunc() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *GatewayConfHandler) ListAllConfigs(ctx context.GContext) error {
	allConfFiles, err := s.service.ListAllConfigs(ctx)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAllConfigs error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, allConfFiles)
}

func (s *GatewayConfHandler) GetConfContent(ctx context.GContext) error {
	param := &GetContentParam{}
	err := echo.QueryParamsBinder(ctx).Strings("name", &param.names).BindError()
	if err != nil {
		common.Log.Warn().Err(err).Str("param", ctx.Request().URL.RawQuery).Msg("Bind query param error.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	if len(param.names) == 0 {
		common.Log.Warn().Str("param", ctx.Request().URL.RawQuery).Msg("Empty query param.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	resp, err := s.service.GetConfContent(ctx, param.names)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAllConfigs error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

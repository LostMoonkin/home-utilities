package handlers

import (
	"encoding/base64"
	"homeserver/common"
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
	group.GET("/listconf", s.ListAllConfs)
	group.GET("/confdetail", s.GetConfContent)
	return nil
}

func (s *GatewayConfHandler) GetAPIPrefix() string {
	return "/api/gateway"
}

func (s *GatewayConfHandler) GetMiddlewareFuncs() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *GatewayConfHandler) ListAllConfs(ctx echo.Context) error {
	allConfFiles, err := s.service.ListAllConfs(ctx)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAllConfs error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, allConfFiles)
}

func (s *GatewayConfHandler) GetConfContent(ctx echo.Context) error {
	param := &GetContentParam{}
	err := echo.QueryParamsBinder(ctx).Strings("name", &param.names).BindError()
	if err != nil {
		common.Log.Warn().Err(err).Str("param", ctx.Request().URL.RawQuery).Msg("Bind query param error.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	contentMap, err := s.service.GetConfContent(ctx, param.names)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAllConfs error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	resp := make(map[string]string)
	for name, content := range contentMap {
		resp[name] = base64.StdEncoding.EncodeToString(content)
	}
	return ctx.JSON(http.StatusOK, resp)
}

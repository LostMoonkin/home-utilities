package handlers

import (
	"homeserver/common"
	"homeserver/context"
	"homeserver/models/gateway"
	"homeserver/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type GatewayConfHandler struct {
	service *services.GatewayConfService
}

func NewGatewayConfHandler(service *services.GatewayConfService) *GatewayConfHandler {
	return &GatewayConfHandler{service: service}
}

func (s *GatewayConfHandler) RegisterRouter(group *echo.Group) error {
	group.GET("/conf/list", wrapHandlerContext(s.ListAll))
	group.GET("/conf", wrapHandlerContext(s.Get))
	group.POST("/conf", wrapHandlerContext(s.Create))
	group.PUT("/conf", wrapHandlerContext(s.Update))
	group.DELETE("/conf", wrapHandlerContext(s.Delete))
	group.POST("/apply", wrapHandlerContext(s.ApplyChange))
	return nil
}

func (s *GatewayConfHandler) GetAPIPrefix() string {
	return "/api/gateway"
}

func (s *GatewayConfHandler) GetMiddlewareFunc() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: 10 * time.Second,
		}),
	}
}

func (s *GatewayConfHandler) Create(ctx context.GContext) error {
	req := &gateway.CreateConfigRequest{}
	if err := ctx.Bind(req); err != nil {
		common.Log.Error().Err(err).Msg("bind request body error")
		return ctx.NoContent(http.StatusBadRequest)
	}
	resp, err := s.service.Create(ctx, req.Name, req.Content)
	if err != nil {
		common.Log.Error().Err(err).Msg("create config error")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *GatewayConfHandler) Update(ctx context.GContext) error {
	req := &gateway.UpdateConfigRequest{}
	if err := ctx.Bind(req); err != nil {
		common.Log.Error().Err(err).Msg("bind request body error")
		return ctx.NoContent(http.StatusBadRequest)
	}
	resp, err := s.service.Update(ctx, req.Name, req.CurrentContent, req.ExpectedContent)
	if err != nil {
		common.Log.Error().Err(err).Msg("update config error")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *GatewayConfHandler) Delete(ctx context.GContext) error {
	return nil
}

func (s *GatewayConfHandler) ApplyChange(ctx context.GContext) error {
	return nil
}

func (s *GatewayConfHandler) ListAll(ctx context.GContext) error {
	allConfFiles, err := s.service.ListAll(ctx)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAll error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, allConfFiles)
}

func (s *GatewayConfHandler) Get(ctx context.GContext) error {
	var nameList []string
	if err := echo.QueryParamsBinder(ctx).Strings("name", &nameList).BindError(); err != nil {
		common.Log.Warn().Err(err).Str("param", ctx.Request().URL.RawQuery).Msg("Bind query param error.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	if len(nameList) == 0 {
		common.Log.Warn().Str("param", ctx.Request().URL.RawQuery).Msg("Empty query param.")
		return ctx.NoContent(http.StatusBadRequest)
	}
	resp, err := s.service.Get(ctx, nameList)
	if err != nil {
		common.Log.Error().Err(err).Msg("ListAll error.")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, resp)
}

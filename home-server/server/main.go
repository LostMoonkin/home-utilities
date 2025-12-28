package main

import (
	"homeserver/common"
	"homeserver/handlers"
	"homeserver/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	common.SetupAppConfig()
	common.InitLogger()

	app := echo.New()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			common.Log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Int64("latencyMS", v.Latency.Abs().Milliseconds()).
				Msg("request")
			return nil
		},
	}))
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	initHandlers(app)
	app.Logger.Fatal(app.Start(common.GetAppConfig().ServerAddress))
}

func initHandlers(app *echo.Echo) {
	// Init handlers
	proxySubHandler := handlers.NewProxySubHandler(services.NewProxySubService())
	handlers.RegisterHandler(proxySubHandler)
	gatewayConfHandler := handlers.NewGatewayConfHandler(services.NewGatewayConfService())
	handlers.RegisterHandler(gatewayConfHandler)

	// Register routers
	handlers.SetupRouter(app)
}

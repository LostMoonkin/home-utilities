package main

import (
	"homeserver/common"
	"homeserver/handlers"
	"homeserver/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	common.SetupAppConfig()
	common.InitLogger()

	app := echo.New()
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

	// Register routers
	handlers.SetupRouter(app)
}

package common

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress        string `env:"SERVER_ADDRESS" envDefault:":8786"`
	LogLevel             string `env:"LOG_LEVEL" envDefault:"DEBUG"`
	EnableFileLogging    bool   `env:"ENABLE_FILE_LOGGING" envDefault:"false"`
	HttpProxy            string `env:"HTTP_PROXY"`
	ClashSubscriptionURL string `ENV:"CLASH_SUB_URL"`
}

var appConfig Config
var initConfigOnce sync.Once

func SetupAppConfig() {
	initConfigOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Sprintf("Get dotenv error: e=%s", err.Error()))
		}

		if err := env.Parse(&appConfig); err != nil {
			panic(fmt.Sprintf("Parse env config error: e=%s", err.Error()))
		}
	})
}

func GetAppConfig() Config {
	return appConfig
}

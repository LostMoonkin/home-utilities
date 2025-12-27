package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

var Log *Logger
var initLogOnce sync.Once

func InitLogger() {
	initLogOnce.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		config := GetAppConfig()

		level, err := zerolog.ParseLevel(config.LogLevel)
		if err != nil {
			panic(fmt.Sprintf("Parse log level error: %s", err.Error()))
		}
		zerolog.SetGlobalLevel(level)

		var writers []io.Writer
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})

		if config.EnableFileLogging {
			fileWrite := lumberjack.Logger{
				Filename:   "logs/server.log",
				MaxSize:    100, // MB
				MaxBackups: 3,
			}
			writers = append(writers, &fileWrite)
		}
		mw := io.MultiWriter(writers...)

		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
		zlogger := zerolog.New(mw).With().Timestamp().Caller().Logger()
		zlogger.Debug().Msg("Logging configured.")
		Log = &Logger{&zlogger}
	})
}

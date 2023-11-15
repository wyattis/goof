package log

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

//go:generate go-enum --marshal --flag

// ENUM(trace, debug, info, warn, error, fatal, panic)
type LogLevel string

var ErrNoOutputs = errors.New("no outputs defined for logger")

type Config struct {
	Level     LogLevel `default:"info"`
	NoConsole bool     `default:"true"`
	Null      bool     `default:"false"`
	File      string   `default:"bbr.log"`
	ErrFile   string   `default:"bbr.err.log"`
}

var Writer io.Writer = os.Stdout
var ErrorWriter io.Writer = os.Stderr
var Logger = zerolog.New(Writer)
var ErrorLogger = zerolog.New(ErrorWriter)
var Level LogLevel = LogLevelInfo

func Init(config Config) (err error) {
	Level = config.Level
	stdout := []io.Writer{}
	stderr := []io.Writer{}
	if !config.NoConsole {
		stdout = append(stdout, os.Stdout)
		stderr = append(stderr, os.Stderr)
	}
	if config.File != "" {
		f, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("Failed to load config.File:\n%w", err)
		}
		stdout = append(stdout, f)
	}
	if config.ErrFile != "" {
		f, err := os.OpenFile(config.ErrFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("Failed to load config.ErrFile:\n%w", err)
		}
		stderr = append(stderr, f)
	}
	if config.Null {
		stdout = append(stdout, io.Discard)
		stderr = append(stderr, io.Discard)
	}
	if len(stdout) == 0 || len(stderr) == 0 {
		return ErrNoOutputs
	}
	Writer = io.MultiWriter(stdout...)
	ErrorWriter = io.MultiWriter(stderr...)
	Logger = zerolog.New(Writer).
		Level(GetLevel(config.Level.String())).
		With().Timestamp().
		Logger()

	ErrorLogger = zerolog.New(ErrorWriter).
		Level(GetLevel(config.Level.String())).
		With().Timestamp().
		Logger()

	Logger.Debug().Msg("Logging initialized")
	return
}

func GetLevel(level string) zerolog.Level {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	return l
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

func Error() *zerolog.Event {
	return ErrorLogger.Error()
}

func Fatal() *zerolog.Event {
	return ErrorLogger.Fatal()
}

func Panic() *zerolog.Event {
	return ErrorLogger.Panic()
}

func Println(v ...any) {
	Logger.Info().Msg(fmt.Sprint(v...))
}

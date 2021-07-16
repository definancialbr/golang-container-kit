package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DevelopmentConfiguration = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	ProductionConfiguration = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "lvl",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
)

type LoggingServiceOption func(*LoggingService)

type LoggingService struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	config        zap.Config
	options       []zap.Option
	name          string
}

func WithName(name string) LoggingServiceOption {
	return func(l *LoggingService) {
		l.name = name
	}
}

func WithOptions(options ...zap.Option) LoggingServiceOption {
	return func(l *LoggingService) {
		l.options = options
	}
}

func WithConfiguration(config zap.Config) LoggingServiceOption {
	return func(l *LoggingService) {
		l.config = config
	}
}

func WithDevelopmentMode() LoggingServiceOption {
	return func(l *LoggingService) {
		l.config = DevelopmentConfiguration
	}
}

func WithProductionMode() LoggingServiceOption {
	return func(l *LoggingService) {
		l.config = ProductionConfiguration
	}
}

func NewLoggingService(options ...LoggingServiceOption) *LoggingService {

	zap.NewProduction()

	l := &LoggingService{
		config: DevelopmentConfiguration,
	}

	for _, option := range options {
		option(l)
	}

	return l

}

func (l *LoggingService) Open() error {

	logger, err := l.config.Build(l.options...)
	if err != nil {
		return err
	}

	if len(l.name) > 0 {
		logger = logger.Named(l.name)
	}

	l.logger = logger
	l.sugaredLogger = logger.Sugar()

	return nil

}

func (l *LoggingService) Close() error {
	// Ignore erros until this issue is not closed  https://github.com/uber-go/zap/issues/880
	l.logger.Sync()
	return nil
}

func (l *LoggingService) Fatal(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *LoggingService) Error(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *LoggingService) Warn(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *LoggingService) Info(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}

func (l *LoggingService) Debug(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *LoggingService) Zap() *zap.Logger {
	return l.logger
}

package boot

import (
	"chatchat/app/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

func Loggersetup() {
	dynamiclevel := zap.NewAtomicLevel() //日志等级

	switch global.Config.Logger.LogLevel {
	case "debug":
		dynamiclevel.SetLevel(zap.DebugLevel)
	case "info":
		dynamiclevel.SetLevel(zap.InfoLevel)
	case "warn":
		dynamiclevel.SetLevel(zap.WarnLevel)
	case "error":
		dynamiclevel.SetLevel(zap.ErrorLevel)

	}
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "message",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       CustomTimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		ConsoleSeparator: "",
	})
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, dynamiclevel),
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(getwritesync()),
			dynamiclevel,
		),
	}
	global.Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(global.Logger)

	global.Logger.Info("initialize logger success")
}

func getwritesync() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Config.Logger.SavePath,
		MaxSize:    global.Config.Logger.MaxSize,
		MaxAge:     global.Config.Logger.MaxBackups,
		MaxBackups: global.Config.Logger.MaxSize,
		LocalTime:  true,
		Compress:   global.Config.Logger.IsCompress,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("[2006-01-02 15:04:05.000]"))
}

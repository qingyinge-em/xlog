package xlog

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(logfile string, to_console bool, levelStr string) (*zap.SugaredLogger, error) {
	level := zap.DebugLevel
	switch strings.ToLower(levelStr) {
	case "err", "error":
		level = zap.ErrorLevel
	case "info":
		level = zap.InfoLevel
	}

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		// return true
		return lev >= level
	})

	encoderConfig := getEncoderConfig()

	getFileLogWriter := func(file string) (writeSyncer zapcore.WriteSyncer) {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   file,
			MaxSize:    10, //MB
			MaxBackups: 5,
			LocalTime:  true,
			// MaxAge:     90, //days
			Compress: false,
		}

		return zapcore.AddSync(lumberJackLogger)
	}

	lowWriteSyncer := getFileLogWriter(logfile)

	highWriteSyncer := getFileLogWriter(func() string {
		errFile := ""
		dotidx := strings.LastIndex(logfile, ".")
		if dotidx >= 0 {
			errFile = logfile[0:dotidx] + "_err" + logfile[dotidx:]
		} else {
			errFile = logfile + "_err"
		}
		return errFile
	}())

	highCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
		highWriteSyncer, highPriority)
	lowCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
		lowWriteSyncer, lowPriority)

	if !to_console {
		return zap.New(zapcore.NewTee(highCore, lowCore), zap.AddCaller()).Sugar(), nil
	}

	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout), lowPriority)
	return zap.New(zapcore.NewTee(highCore, lowCore, consoleCore), zap.AddCaller()).Sugar(), nil
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.LineEnding = "\n"
	encoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
	}
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("0102 15:04:05.000")
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	myTrimmedPath := func(ec zapcore.EntryCaller) string {
		if !ec.Defined {
			return "undefined"
		}

		path := ec.FullPath()

		idx := strings.LastIndexByte(ec.File, '/')
		if idx != -1 {
			idx = strings.LastIndexByte(ec.File[:idx], '/')
			if idx != -1 {
				path = ec.File[idx+1:]
			}
		}

		str := fmt.Sprintf("%s:%d", path, int64(ec.Line))
		return fmt.Sprintf("%-25s", str)
	}
	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(myTrimmedPath(caller))
	}

	encoderConfig.ConsoleSeparator = "  "

	return encoderConfig
}

package log

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	apache_log "github.com/lestrrat-go/apache-logformat"
	"github.com/sirupsen/logrus"
)

const (
	PanicLevel logrus.Level = logrus.PanicLevel
	FatalLevel logrus.Level = logrus.FatalLevel
	ErrorLevel logrus.Level = logrus.ErrorLevel
	WarnLevel  logrus.Level = logrus.WarnLevel
	InfoLevel  logrus.Level = logrus.InfoLevel
	DebugLevel logrus.Level = logrus.DebugLevel
	TraceLevel logrus.Level = logrus.TraceLevel
)

var (
	appLogger    *Logger
	accessLogger *Logger
	config       Targets
)

func InitWithConfiguration(targets Targets) error {
	config = targets
	appLogger = nil
	accessLogger = nil
	return Init()
}

func Init() error {
	if len(config) > 0 {
		for _, target := range config {
			switch len(target.LogTypes) {
			case 0:
				// by default use only app loggers
				configureAppLogger(target)
			case 1:
				switch target.LogTypes[0] {
				case "app":
					configureAppLogger(target)
				case "access":
					configureAccessLogger(target)
				}
			case 2:
				if target.LogTypes[0] != target.LogTypes[1] {
					configureAppLogger(target)
					configureAccessLogger(target)
				}
			}
		}
	} else {
		// No log targets configured: default to stdout for both app and access logs.
		target := Target{
			LogTo:     "stdout",
			LogLevel:  "warning",
			LogFormat: "text",
			LogTypes:  []string{"access", "app"},
		}
		configureAccessLogger(target)
		configureAppLogger(target)
	}
	return nil
}

func configureAppLogger(target Target) {
	logger := logrus.New()
	configureLogger(logger, target)
	if appLogger == nil {
		appLogger = &Logger{}
	}
	aclLogger := &ACLLogger{
		Logger: logger,
	}
	appLogger.loggers = append(appLogger.loggers, aclLogger)
}

func configureAccessLogger(target Target) {
	logger := logrus.New()
	configureLogger(logger, target)
	if accessLogger == nil {
		accessLogger = &Logger{}
	}
	var al *apache_log.ApacheLog
	if target.ACLFormat == "" {
		al, _ = apache_log.New(DefaultApacheLogFormat)
	} else {
		al, _ = apache_log.New(target.ACLFormat)
	}
	if al == nil {
		return
	}
	aclLogger := &ACLLogger{
		Logger:    logger,
		ApacheLog: al,
	}
	accessLogger.loggers = append(accessLogger.loggers, aclLogger)
}

func AppLogger() (*Logger, error) {
	if appLogger == nil {
		if err := Init(); err != nil {
			return nil, err
		}
	}
	if appLogger == nil {
		return nil, errors.New("no application loggers configured")
	}
	return appLogger, nil
}

func AccessLogger() (*Logger, error) {
	if accessLogger == nil {
		if err := Init(); err != nil {
			return nil, err
		}
	}
	if accessLogger == nil {
		return nil, errors.New("no access loggers configured")
	}
	return accessLogger, nil
}

func Log(level logrus.Level, args ...any) {
	if appLogger != nil {
		appLogger.Log(level, args...)
	}
}

func Logf(level logrus.Level, format string, args ...any) {
	if appLogger != nil {
		appLogger.Logf(level, format, args...)
	}
}

func Print(args ...any) {
	if appLogger != nil {
		appLogger.Print(args...)
	}
}

func Trace(args ...any) {
	if appLogger != nil {
		appLogger.Trace(args...)
	}
}

func Debug(args ...any) {
	if appLogger != nil {
		appLogger.Debug(args...)
	}
}

func Info(args ...any) {
	if appLogger != nil {
		appLogger.Info(args...)
	}
}

func Warning(args ...any) {
	if appLogger != nil {
		appLogger.Warning(args...)
	}
}

func Error(args ...any) {
	if appLogger != nil {
		appLogger.Error(args...)
	}
}

func Panic(args ...any) {
	if appLogger != nil {
		appLogger.Panic(args...)
	}
}

func Fatal(args ...any) {
	if appLogger != nil {
		appLogger.Panic(args...)
	}
}

func Printf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Printf(format, args...)
	}
}

func Tracef(format string, args ...any) {
	if appLogger != nil {
		appLogger.Tracef(format, args...)
	}
}

func Debugf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Debugf(format, args...)
	}
}

func Infof(format string, args ...any) {
	if appLogger != nil {
		appLogger.Infof(format, args...)
	}
}

func Warningf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Warningf(format, args...)
	}
}

func Errorf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Errorf(format, args...)
	}
}

func Panicf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Panicf(format, args...)
	}
}

func Fatalf(format string, args ...any) {
	if appLogger != nil {
		appLogger.Fatalf(format, args...)
	}
}

func Fatalln(format string, args ...any) { //nolint:goprintffuncname
	if appLogger != nil {
		appLogger.Fatalln(args...)
	}
}

func WithFields(fields map[string]any, level logrus.Level, args ...any) {
	if appLogger != nil {
		appLogger.WithFields(fields, level, args...)
	}
}

func WithFieldsf(fields map[string]any, level logrus.Level, format string, args ...any) {
	if appLogger != nil {
		appLogger.WithFieldsf(fields, level, format, args...)
	}
}

func configureLogger(logger *logrus.Logger, target Target) {
	switch target.LogFormat {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		})
	case "JSON":
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	switch target.LogTo {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "file":
		dir := filepath.Dir(target.LogFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		logFile, err := os.OpenFile(target.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o666)
		if err != nil {
			Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		logger.SetOutput(logFile)
	case "syslog":
		logger.SetOutput(io.Discard)
		hook, err := NewRFC5424Hook(target)
		if err != nil {
			Warningf("Error configuring Syslog logging: %s", err.Error())
			break
		}
		logger.AddHook(hook)
	}

	switch target.LogLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	}
}

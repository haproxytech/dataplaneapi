package log

import (
	"os"

	apache_log "github.com/lestrrat-go/apache-logformat"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	loggers []*ACLLogger
}

type ACLLogger struct {
	*logrus.Logger
	ApacheLog *apache_log.ApacheLog
}

func (l *Logger) Log(level logrus.Level, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Log(level, args...)
	}
}

func (l *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Logf(level, format, args...)
	}
}

func (l *Logger) Print(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Print(args...)
	}
}

func (l *Logger) Trace(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Trace(args...)
	}
}

func (l *Logger) Debug(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Debug(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Info(args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Warning(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Error(args...)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	for _, log := range l.loggers {
		go log.Panic(args...)
	}
}

func (l *Logger) Printf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Printf(format, args...)
	}
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Tracef(format, args...)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Debugf(format, args...)
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Infof(format, args...)
	}
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Warningf(format, args...)
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Errorf(format, args...)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.Panicf(format, args...)
	}
}

func (l *Logger) Fatalln(args ...interface{}) {
	for _, log := range l.loggers {
		log.Logln(logrus.FatalLevel, args...)
	}
	os.Exit(1)
}

func (l *Logger) Fatal(args ...interface{}) {
	for _, log := range l.loggers {
		log.Log(logrus.FatalLevel, args...)
	}
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	for _, log := range l.loggers {
		log.Logf(logrus.FatalLevel, format, args...)
	}
	os.Exit(1)
}

func (l *Logger) WithFieldsf(fields map[string]interface{}, level logrus.Level, format string, args ...interface{}) {
	for _, log := range l.loggers {
		go log.WithFields(fields).Logf(level, format, args...)
	}
}

func (l *Logger) WithFields(fields map[string]interface{}, level logrus.Level, args ...interface{}) {
	for _, log := range l.loggers {
		go log.WithFields(fields).Log(level, args...)
	}
}

func (l *Logger) Loggers() []*ACLLogger {
	return l.loggers
}

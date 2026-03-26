package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {

	Log.SetOutput(os.Stdout)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Log.SetLevel(logrus.InfoLevel)
		Log.Warnf("Invalid log level '%s', using 'info'", level)
	} else {
		Log.SetLevel(logLevel)
	}

	Log.Infof("Logger initialized with level: %s", Log.GetLevel())
}

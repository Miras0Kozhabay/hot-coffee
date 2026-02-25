package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	Log = logrus.New()
	file, err := os.OpenFile(
		"app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.InfoLevel)

	Log.Info("Application is starting:")
}

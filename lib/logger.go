package lib

import (
	"sync"

	logrus "github.com/sirupsen/logrus"
)

var log *logrus.Logger
var once sync.Once

/*
 TODO:
 1- add logrus formatter logrus.JSONFormatter instance for production env
 2- add file rotate creation hook for logged for production env
*/

func GetLogger() *logrus.Logger {

	/*
		for creating a rotate file
	*/
	// if !debug {
	// 	logLevel = logrus.ErrorLevel
	// }
	// rotateFileHook, err := logrus.NewRotateFileHook(rotatefilehook.RotateFileConfig{
	// 	Filename:   "logs/errors.logrus",
	// 	MaxSize:    50, // megabytes
	// 	MaxBackups: 3,
	// 	MaxAge:     28, //days
	// 	Level:      logLevel,
	// 	Formatter: &logrus.JSONFormatter{
	// 		TimestampFormat: time.RFC822,
	// 	},
	// })

	// if err != nil {
	// 	panic("Failed to initialize file rotate hook: \n, err)
	// }

	// logrus.AddHook(rotateFileHook)

	once.Do(func() {
		var logLevel = logrus.DebugLevel
		logrus.SetLevel(logLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			QuoteEmptyFields:       true,
			TimestampFormat:        "02-01-2006 15:04:05",
		})
		log = logrus.StandardLogger()
	})
	return log
}

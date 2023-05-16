package logging

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
)

type Logger struct {
	*logrus.Entry
}

func (s *Logger) ExtraFields(fields map[string]interface{}) *Logger {
	return &Logger{s.WithFields(fields)}
}

var instance Logger
var once sync.Once

func GetLogger(level string) Logger {
	once.Do(func() {
		logrusLevel, err := logrus.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}
		l := logrus.New()
		l.SetReportCaller(true)
		l.SetOutput(os.Stdout)
		l.SetLevel(logrusLevel)
		instance = Logger{logrus.NewEntry(l)}
	})
	return instance
}

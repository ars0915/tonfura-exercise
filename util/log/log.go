package log

import (
	"github.com/sirupsen/logrus"
)

// Fields type, used to pass to `WithFields`.
type Fields = logrus.Fields

type formatter struct {
	logrus.JSONFormatter
}

func (f formatter) Format(e *logrus.Entry) ([]byte, error) {
	f.JSONFormatter.TimestampFormat = "2006/01/02 15:04:05"
	e.Time = e.Time.UTC()
	return f.JSONFormatter.Format(e)
}

func init() {
	logrus.SetFormatter(&formatter{})
}

// SetLevel set log level
func SetLogLevel(level string) error {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	logrus.SetLevel(lv)
	return nil
}

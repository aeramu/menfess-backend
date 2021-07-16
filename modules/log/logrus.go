package log

import (
	"github.com/aeramu/menfess-backend/service"
	loglib "github.com/sirupsen/logrus"
)

func NewLogModule() service.LogModule {
	return &logModule{}
}

type logModule struct {}

func (m *logModule) Log(err error, payload interface{}, message string) {
	loglib.WithFields(loglib.Fields{
		"err": err,
		"payload": payload,
	}).Errorln(message)
}

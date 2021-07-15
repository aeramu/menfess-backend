package log

import loglib "github.com/sirupsen/logrus"

type logModule struct {}

func (m *logModule) Log(err error, payload interface{}, message string) {
	loglib.WithFields(loglib.Fields{
		"err": err,
		"payload": payload,
	}).Errorln(message)
}

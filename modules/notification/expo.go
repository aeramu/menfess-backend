package notification

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	expoURL = "https://exp.host/--/api/v2/push/send"
	contentType = "application/json"
)

func (m *notificationModule) sendNotification(tokens []string, title, body string, data interface{}) error {
	req := notificationRequest{
		To:    tokens,
		Title: title,
		Body:  body,
		Data:  data,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	go func() {
		_, err := http.Post(expoURL, contentType, bytes.NewReader(b))
		if err != nil {
			logrus.Errorln("[SendNotification] Failed send http request:", err)
		}
	}()

	return nil
}

type notificationRequest struct {
	To    []string    `json:"to"`
	Title string      `json:"title"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data"`
}

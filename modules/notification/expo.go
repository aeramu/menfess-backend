package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const (
	expoURL = "https://exp.host/--/api/v2/push/send"
	contentType = "application/json"
)

func (m *notificationModule) sendNotification(ctx context.Context, tokens []string, title, body string, data interface{}) error {
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
		httpResp, err := http.Post(expoURL, contentType, bytes.NewReader(b))
		if err != nil {
			logrus.Errorln("[SendNotification] Failed send http request:", err)
			return
		}
		b, err = ioutil.ReadAll(httpResp.Body)
		if err != nil {
			logrus.Errorln("[SendNotification] Failed read http response:", err)
			return
		}
		var resp notificationResponse
		if err := json.Unmarshal(b, &resp); err != nil {
			logrus.Errorln("[SendNotification] Failed unmarshall response:", err)
			return
		}
		m.processNotificationResponse(ctx, resp, tokens)
	}()

	return nil
}

func (m *notificationModule) processNotificationResponse(ctx context.Context, resp notificationResponse, tokens []string) {
	for i, v := range resp.Data {
		if v.Status == "ok" {
			break
		}
		if v.Details.Error == "DeviceNotRegistered" {
			err := m.pushToken.Query().Equal("token", tokens[i]).DeleteMany(ctx)
			if err != nil {
				logrus.Errorln("Failed delete invalid expo token, err:", err)
			}
		}
	}
}

type notificationRequest struct {
	To    []string    `json:"to"`
	Title string      `json:"title"`
	Body  string      `json:"body"`
	Data  interface{} `json:"data"`
}

type notificationResponse struct {
	Data []pushTicket `json:"data"`
}

type pushTicket struct {
	ID      string           `json:"id"`
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Details pushTicketDetail `json:"details"`
}

type pushTicketDetail struct {
	Error string `json:"error"`
	Fault string `json:"fault"`
}
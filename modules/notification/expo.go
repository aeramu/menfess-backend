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

func (m *notificationModule) sendNotification(_ context.Context, tokens []string, title, body string, data interface{}) error {
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
		m.processNotificationResponse(context.Background(), resp, tokens)
	}()

	return nil
}

func (m *notificationModule) processNotificationResponse(ctx context.Context, resp notificationResponse, tokens []string) {
	var errorToken []string
	for i, v := range resp.Data {
		if v.Status == "ok" {
			continue
		}
		if v.Details.Error == "DeviceNotRegistered" {
			errorToken = append(errorToken, tokens[i])
			continue
		}
		logrus.Errorln("Expo token error, token:", tokens[i], v.Details.Error)
	}
	if len(errorToken) > 0 {
		err := m.pushToken.Query().In("token", errorToken).Delete(ctx)
		if err != nil {
			logrus.Errorln("Failed delete invalid expo token, err:", err)
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
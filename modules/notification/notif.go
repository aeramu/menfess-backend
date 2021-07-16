package notification

import (
	"context"
	"fmt"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
)

func NewNotificationModule() service.NotificationModule {
	return &notificationModule{}
}

type notificationModule struct {

}

func (m *notificationModule) AddPushToken(ctx context.Context, userID string, pushToken string) error {
	fmt.Printf("add push token [%s] to user: %s\n", pushToken, userID)
	return nil
}

func (m *notificationModule) SendLikeNotification(ctx context.Context, user entity.User, post entity.Post) error {
	fmt.Printf("send notification [%s like your post] to %s\n", user.Profile.Name, post.User.ID)
	return nil
}


package notification

import (
	"context"
	"fmt"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
	"github.com/aeramu/menfess-backend/utils"
	"github.com/aeramu/mongolib"
)

func NewNotificationModule(db *mongolib.Database) service.NotificationModule {
	return &notificationModule{
		pushToken: db.Coll("push_token"),
	}
}

type notificationModule struct {
	pushToken    *mongolib.Collection
}

const(
	likeNotificationTitle = "%s like your post"
	commentNotificationTitle = "%s comment on your post"
	newPostNotificationTitle = "Someone post a menfess just now"
)

func (m *notificationModule) AddPushToken(ctx context.Context, userID string, pushToken string) error {
	if err := m.insertPushToken(ctx, userID, pushToken); err != nil {
		return err
	}
	return nil
}

func (m *notificationModule) RemovePushToken(ctx context.Context, userID string, pushToken string) error {
	if err := m.removePushToken(ctx, userID, pushToken); err != nil {
		if err == mongolib.ErrNotFound {
			return constants.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (m *notificationModule) SendLikeNotification(ctx context.Context, user entity.User, post entity.Post) error {
	tokens, err := m.findPushToken(ctx, post.User.ID)
	if err != nil {
		return err
	}

	if err := m.sendNotification(
		ctx,
		tokens,
		fmt.Sprintf(likeNotificationTitle, user.Profile.Name),
		post.Body,
		Data{PostID: post.ID},
	); err != nil {
		return err
	}

	return nil
}

func (m *notificationModule) SendCommentNotification(ctx context.Context, comment entity.Post, parent entity.Post) error {
	tokens, err := m.findPushToken(ctx, parent.User.ID)
	if err != nil {
		return err
	}

	if err := m.sendNotification(
		ctx,
		tokens,
		fmt.Sprintf(commentNotificationTitle, comment.User.Profile.Name),
		comment.Body,
		Data{PostID: parent.ID},
	); err != nil {
		return err
	}

	return nil
}

func (m *notificationModule) BroadcastNewPostNotification(ctx context.Context, post entity.Post) error {
	pushTokens, err := m.findAllPushToken(ctx)
	if err != nil {
		return err
	}

	var tokens []string
	for _, v := range pushTokens {
		if v.UserID.Hex() != post.User.ID {
			tokens = append(tokens, v.Token)
		}
	}

	if !utils.RandomChance(100, 100) {
		return nil
	}

	if err := m.sendNotification(
		ctx,
		tokens,
		fmt.Sprintf(newPostNotificationTitle),
		post.Body,
		Data{PostID: post.ID},
	); err != nil {
		return err
	}

	return nil
}

type Data struct {
	PostID string `json:"postID"`
}
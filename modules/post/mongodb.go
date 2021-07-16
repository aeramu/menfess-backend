package post

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service/api"
	"github.com/aeramu/mongolib"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type postModule struct {
	post *mongolib.Collection
}

func (m *postModule) FindPostByID(ctx context.Context, id string, userID string) (*entity.Post, error) {
	var model AggregatePostList
	if err := m.post.Aggregate().
		Match(mongolib.Filter{}.Equal("_id", mongolib.ObjectID(id))).
		Lookup("user", "author_id", "_id", "author").
		Lookup("user", "user_id", "_id", "user").
		Exec(ctx).Consume(&model);
	err != nil {
		return nil, err
	}
	if len(model) < 1 {
		return nil, constants.ErrPostNotFound
	}
	return model[0].Entity(userID), nil
}

func (m *postModule) FindPostListByParentIDAndAuthorIDs(ctx context.Context, parentID string, authorIDs []string, userID string, pagination api.PaginationReq) ([]entity.Post, *api.PaginationRes, error) {
	var model AggregatePostList
	if err := m.post.Aggregate().
		Match(mongolib.Filter{}.
			Equal("parent_id", mongolib.ObjectID(parentID)).
			// TODO: Need fixing
			In("author_id", authorIDs).
			LessThan("_id", mongolib.ObjectID(pagination.After))).
		Sort("_id", mongolib.Descending).
		Limit(pagination.First).
		Lookup("user", "author_id", "_id", "author").
		Lookup("user", "user_id", "_id", "user").
		Exec(ctx).Consume(&model);
		err != nil {
		return nil, nil, err
	}

	return model.Entity(userID), &api.PaginationRes{
		EndCursor:   model[len(model)-1].ID.Hex(),
		HasNextPage: len(model) >= pagination.First,
	}, nil
}

func (m *postModule) SavePost(ctx context.Context, post entity.Post) error {
	id := mongolib.NewObjectID()
	if err := m.post.Save(ctx, id, Post{
		ID:           id,
		Body:         post.Body,
		RepliesCount: post.RepliesCount,
		Likes:        map[string]bool{},
		ParentID:     mongolib.ObjectID(post.Parent.ID),
		AuthorID:     mongolib.ObjectID(post.Author.ID),
		User:         mongolib.ObjectID(post.User.ID),
	}); err != nil {
		return err
	}

	return nil
}

func (m *postModule) LikePost(ctx context.Context, postID string, userID string) error {
	panic("implement me")
}

func (m *postModule) UnlikePost(ctx context.Context, postID string, userID string) error {
	panic("implement me")
}

type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	Body         string             `bson:"body"`
	RepliesCount int                `bson:"replies_count"`
	Likes        map[string]bool    `bson:"likes"`
	ParentID     primitive.ObjectID `bson:"parent_id"`
	AuthorID     primitive.ObjectID `bson:"author_id"`
	User         primitive.ObjectID `bson:"user_id"`
}

type AggregatePost struct {
	ID           primitive.ObjectID `bson:"_id"`
	Body         string             `bson:"body"`
	RepliesCount int                `bson:"replies_count"`
	Likes        map[string]bool    `bson:"likes"`
	ParentID     primitive.ObjectID `bson:"parent_id"`
	Author       User               `bson:"author"`
	User         User               `bson:"user"`
}

type User struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Avatar string             `bson:"avatar"`
}

func (p AggregatePost) Entity(userID string) *entity.Post {
	isLiked := false
	if _, ok := p.Likes[userID]; ok {
		isLiked = true
	}
	return &entity.Post{
		ID:           p.ID.Hex(),
		Body:         p.Body,
		Timestamp:    p.ID.Timestamp().Unix(),
		RepliesCount: p.RepliesCount,
		LikesCount:   len(p.Likes),
		IsLiked:      isLiked,
		Parent:       &entity.Post{ID: p.ParentID.Hex()},
		Author:       entity.User{
			ID:      p.Author.ID.Hex(),
			Profile: entity.Profile{
				Name:   p.Author.Name,
				Avatar: p.Author.Avatar,
			},
		},
		User:         entity.User{
			ID:      p.User.ID.Hex(),
			Profile: entity.Profile{
				Name:   p.User.Name,
				Avatar: p.User.Avatar,
			},
		},
	}
}

type AggregatePostList []AggregatePost

func (p AggregatePostList) Entity(userID string) []entity.Post {
	var entities []entity.Post
	for _, v := range p {
		entities = append(entities, *v.Entity(userID))
	}
	return entities
}
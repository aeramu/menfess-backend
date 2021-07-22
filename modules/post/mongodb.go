package post

import (
	"context"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
	"github.com/aeramu/menfess-backend/service/api"
	"github.com/aeramu/mongolib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewPostModule(db *mongolib.Database) service.PostModule {
	return &postModule{post: db.Coll("post")}
}

type postModule struct {
	post *mongolib.Collection
}

func (m *postModule) FindPostByID(ctx context.Context, id string, userID string) (*entity.Post, error) {
	var model AggregatePostList
	if err := m.post.Aggregate().
		Match(mongolib.Filter().Equal("_id", mongolib.ObjectID(id))).
		Lookup("user", "author_id", "_id", "author").
		Unwind("$author").
		Lookup("user", "user_id", "_id", "user").
		Unwind("$user").
		Lookup("post", "_id", "parent_id", "replies").
		AddField("replies_count", bson.D{{"$size", "$replies"}}).
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
	if pagination.After == "" {
		pagination.After = "ffffffffffffffffffffffff"
	}
	if err := m.post.Aggregate().
		Match(mongolib.Filter().
			Equal("parent_id", mongolib.ObjectID(parentID)).
			// TODO: Need fixing filter author ids
			//In("author_id", authorIDs).
			LessThan("_id", mongolib.ObjectID(pagination.After))).
		Sort("_id", mongolib.Descending).
		Limit(pagination.First).
		Lookup("user", "author_id", "_id", "author").
		Unwind("$author").
		Lookup("user", "user_id", "_id", "user").
		Unwind("$user").
		Lookup("post", "_id", "parent_id", "replies").
		AddField("replies_count", bson.D{{"$size", "$replies"}}).
		Exec(ctx).Consume(&model);
		err != nil {
		return nil, nil, err
	}

	endCursor := ""
	if len(model) > 0 {
		endCursor = model[len(model)-1].ID.Hex()
	}
	return model.Entity(userID), &api.PaginationRes{
		EndCursor: endCursor,
		HasNextPage: len(model) >= pagination.First,
	}, nil
}

func (m *postModule) InsertPost(ctx context.Context, post entity.Post) (string, error) {
	id := mongolib.NewObjectID()
	if err := m.post.Save(ctx, id, Post{
		ID:           id,
		Body:         post.Body,
		Likes:        map[string]bool{},
		ParentID:     mongolib.ObjectID(post.Parent.ID),
		AuthorID:     mongolib.ObjectID(post.Author.ID),
		User:         mongolib.ObjectID(post.User.ID),
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (m *postModule) LikePost(ctx context.Context, postID string, userID string) error {
	var model Post
	if err := m.post.Query().
		Equal("_id", mongolib.ObjectID(postID)).
		FindOne(ctx).Consume(&model);
	err != nil {
		return err
	}

	model.Likes[userID] = true

	if err := m.post.Save(ctx, mongolib.ObjectID(postID), model); err != nil {
		return err
	}

	return nil
}

func (m *postModule) UnlikePost(ctx context.Context, postID string, userID string) error {
	var model Post
	if err := m.post.Query().
		Equal("_id", mongolib.ObjectID(postID)).
		FindOne(ctx).Consume(&model);
		err != nil {
		return err
	}

	delete(model.Likes, userID)

	if err := m.post.Save(ctx, mongolib.ObjectID(postID), model); err != nil {
		return err
	}

	return nil
}

type Post struct {
	ID           primitive.ObjectID `bson:"_id"`
	Body         string             `bson:"body"`
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
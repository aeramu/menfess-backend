package graphql

import (
	"context"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/graph-gophers/graphql-go"
)

type User struct {
	ID     graphql.ID
	Name   string
	Avatar string
	Bio    string
}

type UserEdge struct {
	Node   User
	Cursor graphql.ID
}

type UserConnection struct {
	Edges    []UserEdge
	PageInfo PageInfo
}

func ResolveUser(user entity.User) User {
	return User{
		ID:     graphql.ID(user.ID),
		Name:   user.Profile.Name,
		Avatar: user.Profile.Avatar,
		Bio:    user.Profile.Bio,
	}
}

func ResolveUserEdges(posts []entity.User) []UserEdge {
	var edges []UserEdge
	for _, v := range posts {
		edges = append(edges, UserEdge{
			Node:   ResolveUser(v),
			Cursor: graphql.ID(v.ID),
		})
	}
	return edges
}

type Post struct {
	ID           graphql.ID
	Body         string
	Timestamp    int32
	Author       User
	LikesCount   int32
	RepliesCount int32
	IsLiked      bool
}

func (p Post) Replies(ctx context.Context, input struct{
	First int
	After *graphql.ID
}) PostConnection {
	return PostConnection{}
}

type PostEdge struct {
	Node   Post
	Cursor graphql.ID
}

type PostConnection struct {
	Edges    []PostEdge
	PageInfo PageInfo
}

func ResolvePost(post entity.Post) Post {
	return Post{
		ID:           graphql.ID(post.ID),
		Body:         post.Body,
		Timestamp:    int32(post.Timestamp),
		Author:       ResolveUser(post.Author),
		LikesCount:   int32(post.LikesCount),
		RepliesCount: int32(post.RepliesCount),
		IsLiked:      post.IsLiked,
	}
}

func ResolvePostEdges(posts []entity.Post) []PostEdge {
	var edges []PostEdge
	for _, v := range posts {
		edges = append(edges, PostEdge{
			Node:   ResolvePost(v),
			Cursor: graphql.ID(v.ID),
		})
	}
	return edges
}

type PageInfo struct {
	EndCursor   graphql.ID
	HasNextPage bool
}
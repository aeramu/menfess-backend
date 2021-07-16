package graphql

import (
	"context"
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

type PageInfo struct {
	EndCursor   graphql.ID
	HasNextPage bool
}
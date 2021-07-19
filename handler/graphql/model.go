package graphql

import (
	"context"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service/api"
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
	*Resolver
	ID           graphql.ID
	Body         string
	Timestamp    int32
	Author       User
	LikesCount   int32
	RepliesCount int32
	IsLiked      bool
}

func (p Post) Replies(ctx context.Context, input struct{
	First int32
	After *graphql.ID
}) PostConnection {
	token, err := DecodeToken(ctx)
	if err != nil {
		return PostConnection{}
	}
	req := api.GetPostListReq{
		UserID:     token.UserID,
		ParentID:   string(p.ID),
		Pagination: api.PaginationReq{
			First: int(input.First),
		},
	}
	if input.After != nil {
		req.Pagination.After = string(*input.After)
	}
	res, err := p.svc.GetPostList(ctx, req)
	if err != nil {
		return PostConnection{}
	}

	return PostConnection{
		Edges:    ResolvePostEdges(p.Resolver, res.PostList),
		PageInfo: PageInfo{
			EndCursor:   graphql.ID(res.Pagination.EndCursor),
			HasNextPage: res.Pagination.HasNextPage,
		},
	}
}

type PostEdge struct {
	Node   Post
	Cursor graphql.ID
}

type PostConnection struct {
	Edges    []PostEdge
	PageInfo PageInfo
}

func ResolvePost(r *Resolver, post entity.Post) Post {
	return Post{
		Resolver:     r,
		ID:           graphql.ID(post.ID),
		Body:         post.Body,
		Timestamp:    int32(post.Timestamp),
		Author:       ResolveUser(post.Author),
		LikesCount:   int32(post.LikesCount),
		RepliesCount: int32(post.RepliesCount),
		IsLiked:      post.IsLiked,
	}
}

func ResolvePostEdges(r *Resolver, posts []entity.Post) []PostEdge {
	var edges []PostEdge
	for _, v := range posts {
		edges = append(edges, PostEdge{
			Node:   ResolvePost(r, v),
			Cursor: graphql.ID(v.ID),
		})
	}
	return edges
}

type PageInfo struct {
	EndCursor   graphql.ID
	HasNextPage bool
}
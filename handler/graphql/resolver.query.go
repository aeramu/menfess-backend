package graphql

import (
	"context"
	"github.com/aeramu/menfess-backend/service/api"
	"github.com/graph-gophers/graphql-go"
)

func (r *Resolver) Post(ctx context.Context, input struct{
	ID graphql.ID
}) PostResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return PostResponse{
			Error:   Error(err),
		}
	}
	res, err := r.svc.GetPost(ctx, api.GetPostReq{
		ID: string(input.ID),
		UserID: token.UserID,
	})
	if err != nil {
		return PostResponse{Error: Error(err)}
	}
	return PostResponse{
		Payload: Post{
			ID:           graphql.ID(res.Post.ID),
			Body:         res.Post.Body,
			Timestamp:    int32(res.Post.Timestamp),
			Author:       User{
				ID:     graphql.ID(res.Post.Author.ID),
				Name:   res.Post.Author.Profile.Name,
				Avatar: res.Post.Author.Profile.Avatar,
				Bio:    res.Post.Author.Profile.Bio,
			},
			LikesCount:   int32(res.Post.LikesCount),
			RepliesCount: int32(res.Post.RepliesCount),
			IsLiked:      res.Post.IsLiked,
		},
		Error:   NoError,
	}
}

func (r *Resolver) Posts(ctx context.Context, input struct{
	First int
	After *graphql.ID
	Filter *[]graphql.ID
}) PostsResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return PostsResponse{
			Error:   Error(err),
		}
	}
	res, err := r.svc.GetPostList(ctx, api.GetPostListReq{
		ParentID:   "",
		AuthorIDs:  nil,
		UserID:     token.UserID,
		Pagination: api.PaginationReq{
			First: input.First,
			After: "",
		},
	})
	if err != nil {
		return PostsResponse{Error: Error(err)}
	}
	return PostsResponse{
		Payload: PostConnection{
			Edges:    nil,
			PageInfo: PageInfo{
				EndCursor:   graphql.ID(res.Pagination.EndCursor),
				HasNextPage: res.Pagination.HasNextPage,
			},
		},
		Error:   NoError,
	}
}

func (r *Resolver) Me(ctx context.Context) MeResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return MeResponse{
			Error:   Error(err),
		}
	}
	res, err := r.svc.GetUser(ctx, api.GetUserReq{
		ID: token.UserID,
	})
	if err != nil {
		return MeResponse{Error: Error(err)}
	}
	return MeResponse{
		Payload: User{
			ID:     graphql.ID(res.User.ID),
			Name:   res.User.Profile.Name,
			Avatar: res.User.Profile.Avatar,
			Bio:    res.User.Profile.Bio,
		},
		Error:   NoError,
	}
}

func (r *Resolver) Menfess(ctx context.Context) MenfessResponse {
	_, err := r.svc.GetMenfessList(ctx, api.GetMenfessListReq{})
	if err != nil {
		return MenfessResponse{Error: Error(err)}
	}
	return MenfessResponse{
		Payload: UserConnection{},
		Error:   NoError,
	}
}
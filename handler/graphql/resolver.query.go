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
		Payload: ResolvePost(r, res.Post),
		Error:   NoError,
	}
}

func (r *Resolver) Feed(ctx context.Context, input struct{
	First int32
	After *graphql.ID
	Filter *[]graphql.ID
}) FeedResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return FeedResponse{
			Error:   Error(err),
		}
	}
	req := api.GetPostListReq{
		UserID:     token.UserID,
		Pagination: api.PaginationReq{
			First: int(input.First),
		},
	}
	if input.After != nil {
		req.Pagination.After = string(*input.After)
	}
	res, err := r.svc.GetPostList(ctx, req)
	if err != nil {
		return FeedResponse{Error: Error(err)}
	}
	return FeedResponse{
		Payload: PostConnection{
			Edges:    ResolvePostEdges(r, res.PostList),
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
	res, err := r.svc.GetMenfessList(ctx, api.GetMenfessListReq{})
	if err != nil {
		return MenfessResponse{Error: Error(err)}
	}
	return MenfessResponse{
		Payload: UserConnection{
			Edges: ResolveUserEdges(res.MenfessList),
		},
		Error:   NoError,
	}
}

func (r *Resolver) Avatars(ctx context.Context) AvatarsResponse {
	return AvatarsResponse{
		Payload: []string{},
		Error:   NoError,
	}
}
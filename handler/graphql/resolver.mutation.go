package graphql

import (
	"context"
	"github.com/aeramu/menfess-backend/service/api"
	"github.com/graph-gophers/graphql-go"
)

func (r *Resolver) Login(ctx context.Context, input struct{
	Email string
	Password string
	PushToken string
}) AuthResponse {
	res, err := r.svc.Login(ctx, api.LoginReq{
		Email:     input.Email,
		Password:  input.Password,
		PushToken: input.PushToken,
	})
	if err != nil {
		return AuthResponse{Error: Error(err)}
	}
	return AuthResponse{
		Token: res.Token,
		Error: NoError,
	}
}

func (r *Resolver) Logout(ctx context.Context, input struct{
	PushToken string
}) BasicMutationResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return BasicMutationResponse{
			Error:   Error(err),
		}
	}
	_, err = r.svc.Logout(ctx, api.LogoutReq{
		UserID:    token.UserID,
		PushToken: input.PushToken,
	})
	if err != nil {
		return BasicMutationResponse{Error: Error(err)}
	}
	return BasicMutationResponse{
		Error: NoError,
	}
}

func (r *Resolver) Register(ctx context.Context, input struct{
	PushToken string
}) AuthResponse {
	res, err := r.svc.Register(ctx, api.RegisterReq{
		PushToken: input.PushToken,
	})
	if err != nil {
		return AuthResponse{Error: Error(err)}
	}
	return AuthResponse{
		Token: res.Token,
		Error: NoError,
	}
}

func (r *Resolver) UpdateProfile(ctx context.Context, input struct{
	Name string
	Avatar string
	Bio string
}) BasicMutationResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return BasicMutationResponse{
			Error:   Error(err),
		}
	}
	res, err := r.svc.UpdateProfile(ctx, api.UpdateProfileReq{
		ID:     token.UserID,
		Name:   input.Name,
		Avatar: input.Avatar,
		Bio:    input.Bio,
	})
	if err != nil {
		return BasicMutationResponse{Error: Error(err)}
	}
	return BasicMutationResponse{
		Message: res.Message,
		Error:   NoError,
	}
}

func (r *Resolver) CreatePost(ctx context.Context, input struct{
	Body string
	AuthorID *graphql.ID
	ParentID *graphql.ID
}) BasicMutationResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return BasicMutationResponse{
			Error:   Error(err),
		}
	}

	req := api.CreatePostReq{
		Body:     input.Body,
		UserID:   token.UserID,
	}
	if input.AuthorID ==  nil {
		req.AuthorID = token.UserID
	} else {
		req.AuthorID = string(*input.AuthorID)
	}
	if input.ParentID != nil {
		req.ParentID = string(*input.ParentID)
	}
	res, err := r.svc.CreatePost(ctx, req)
	if err != nil {
		return BasicMutationResponse{Error: Error(err)}
	}
	return BasicMutationResponse{
		Message: res.Message,
		Error:   NoError,
	}
}

func (r *Resolver) LikePost(ctx context.Context, input struct{
	ID graphql.ID
}) BasicMutationResponse {
	token, err := DecodeToken(ctx)
	if err != nil {
		return BasicMutationResponse{
			Error:   Error(err),
		}
	}
	res, err := r.svc.LikePost(ctx, api.LikePostReq{
		PostID: string(input.ID),
		UserID: token.UserID,
	})
	if err != nil {
		return BasicMutationResponse{Error: Error(err)}
	}
	return BasicMutationResponse{
		Message: res.Message,
		Error:   NoError,
	}
}


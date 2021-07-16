package graphql

import (
	"context"
	"errors"
	"github.com/aeramu/menfess-backend/constants"
	"github.com/aeramu/menfess-backend/modules/auth"
)

type Token struct {
	UserID string
}

func DecodeToken(ctx context.Context) (*Token, error) {
	tokenString, ok := ctx.Value(constants.AuthorizationKey).(string)
	if !ok {
		return nil, errors.New("token is required")
	}
	claim, err := auth.DecodeToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	return &Token{
		UserID: claim.UserID,
	}, nil
}

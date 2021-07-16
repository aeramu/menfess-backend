package auth

import (
	"context"
	"errors"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type authModule struct {

}

type Claims struct {
	UserID string
	jwt.StandardClaims
}

var key = []byte("menfess")

func (m *authModule) GenerateToken(ctx context.Context, user entity.User) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: user.ID,
	}).SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func DecodeToken(ctx context.Context, tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

func (m *authModule) ComparePassword(ctx context.Context, hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (m *authModule) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}


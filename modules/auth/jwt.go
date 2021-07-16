package auth

import (
	"context"
	"errors"
	"github.com/aeramu/menfess-backend/entity"
	"github.com/aeramu/menfess-backend/service"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthModule() service.AuthModule {
	return &AuthModule{}
}

type AuthModule struct {

}

type Claims struct {
	UserID string
	jwt.StandardClaims
}

var key = []byte("menfess")

func (m *AuthModule) GenerateToken(ctx context.Context, user entity.User) (string, error) {
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

func (m *AuthModule) ComparePassword(ctx context.Context, hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (m *AuthModule) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}


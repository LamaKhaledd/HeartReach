package services

import (
	"context"
	"errors"

	"github.com/LamaKhaledd/HeartReach/internal/db"
	"github.com/LamaKhaledd/HeartReach/internal/utils"
)

type LoginService struct {
	Queries *db.Queries
	JwtKey  string
}

func (l *LoginService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := l.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !utils.IsPasswordMatches(password, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := utils.CreateJwtAccessToken(l.JwtKey, user)
	if err != nil {
		return "", err
	}

	return token, nil
}

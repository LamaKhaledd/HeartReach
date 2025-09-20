package services

import (
	"context"
	"database/sql"

	"github.com/LamaKhaledd/HeartReach/internal/db"
	"github.com/LamaKhaledd/HeartReach/internal/utils"
)

type RegisterService struct {
	Queries *db.Queries
	JwtKey  string
}

func (r *RegisterService) Register(ctx context.Context, email, userName, password, phoneNumber, role, location string) (string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := r.Queries.CreateUser(ctx, db.CreateUserParams{
		Email:       email,
		UserName:    userName,
		Password:    hashedPassword,
		PhoneNumber: sql.NullString{String: phoneNumber, Valid: phoneNumber != ""},
		Role:        role,
		Location:    sql.NullString{String: location, Valid: location != ""},
	})
	if err != nil {
		return "", err
	}

	token, err := utils.CreateJwtAccessToken(r.JwtKey, user)
	if err != nil {
		return "", err
	}

	return token, nil
}

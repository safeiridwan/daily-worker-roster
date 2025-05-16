package user

import (
	"backend/internal/service/user"
	"context"
	"database/sql"
	"github.com/rs/zerolog"
	"time"
)

type Usecase interface {
	CreateUser(ctx context.Context, input UserIn) (*CreateUserOut, error)
}

type usecase struct {
	userService user.Service
}

func NewUsecase(userService user.Service) Usecase {
	return &usecase{
		userService: userService,
	}
}

func (u *usecase) CreateUser(ctx context.Context, input UserIn) (*CreateUserOut, error) {
	exist, err := u.userService.GetByEmail(ctx, input.Email)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	if exist != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("User already exists")
		return nil, err
	}

	plainPassword, hashPassword, err := GeneratePassword()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	userEntity := user.User{
		Email:     sql.NullString{String: input.Email, Valid: true},
		Name:      sql.NullString{String: input.Name, Valid: true},
		Role:      sql.NullString{String: input.Role, Valid: true},
		Password:  sql.NullString{String: hashPassword, Valid: true},
		CreatedBy: sql.NullString{String: "System", Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := u.userService.SaveUser(ctx, userEntity); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &CreateUserOut{
		Password: plainPassword,
	}

	return out, nil
}

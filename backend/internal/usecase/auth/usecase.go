package auth

import (
	"backend/internal/config"
	"backend/internal/service/middleware"
	"backend/internal/service/user"
	"context"
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Usecase interface {
	Login(ctx context.Context, input LoginIn) (*LoginOut, error)
	UserInfo(ctx context.Context) (*InfoOut, error)
}

type usecase struct {
	service user.Service
	ja      *jwtauth.JWTAuth
	cfg     *config.Config
}

func NewUsecase(
	service user.Service,
	ja *jwtauth.JWTAuth,
	cfg *config.Config,
) Usecase {
	return &usecase{
		service: service,
		ja:      ja,
		cfg:     cfg,
	}
}

func (u *usecase) Login(ctx context.Context, input LoginIn) (*LoginOut, error) {
	userEntity, err := u.service.GetByEmail(ctx, input.Username)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	if userEntity == nil {
		err := errors.New("username or password is incorrect")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(input.Password)); err != nil {
		err = errors.New("invalid password")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	_, tokenstring, err := u.ja.Encode(map[string]interface{}{
		"user_uid": userEntity.UID.String(),
		"exp":      time.Now().Add(time.Duration(u.cfg.JWTExpiration) * time.Second).Unix(),
	})

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &LoginOut{}
	out.Token = tokenstring

	return out, nil
}

func (u *usecase) UserInfo(ctx context.Context) (*InfoOut, error) {
	identity := middleware.CurrentUser(ctx)

	userData, err := u.service.GetByUID(ctx, identity.UserUID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	if userData == nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("user not found")
		return nil, err
	}

	res := &InfoOut{}
	res.Populate(userData)

	return res, nil
}

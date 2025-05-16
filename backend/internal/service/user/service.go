package user

import (
	"context"
	"github.com/rs/zerolog"
)

type Service interface {
	GetByEmail(ctx context.Context, email string) (*UserOut, error)
	SaveUser(ctx context.Context, user User) error
}

type service struct {
	repo UserRepository
}

func NewService(repo UserRepository) Service {
	return &service{repo}
}

func (s *service) GetByEmail(ctx context.Context, email string) (*UserOut, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	if user == nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("User not found")
		return nil, err
	}

	out := &UserOut{}
	out.PopulateFromEntity(user)

	return out, nil
}

func (s *service) SaveUser(ctx context.Context, user User) error {
	err := s.repo.SaveUser(ctx, user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

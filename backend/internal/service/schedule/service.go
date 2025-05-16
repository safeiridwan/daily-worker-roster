package schedule

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type Service interface {
	SaveSchedule(ctx context.Context, schedule Schedule) error
	GetByRangeDate(ctx context.Context, date, startTime, endTime time.Time) (*SchedulesOut, error)
}

type service struct {
	repo ScheduleRepository
}

func NewService(repo ScheduleRepository) Service {
	return &service{repo}
}

func (s *service) SaveSchedule(ctx context.Context, schedule Schedule) error {
	err := s.repo.SaveSchedule(ctx, schedule)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s *service) GetByRangeDate(ctx context.Context, date, startTime, endTime time.Time) (*SchedulesOut, error) {
	schedules, total, err := s.repo.GetByRangeDate(ctx, date, startTime, endTime)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &SchedulesOut{}
	out.PopulateFromEntity(schedules, total)

	return out, nil
}

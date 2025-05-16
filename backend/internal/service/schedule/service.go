package schedule

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

type Service interface {
	SaveSchedule(ctx context.Context, schedule Schedule) error
	UpdateSchedule(ctx context.Context, uid string, schedule Schedule) error
	GetByRangeDate(ctx context.Context, date, startTime, endTime time.Time) (*SchedulesOut, error)
	GetByUID(ctx context.Context, uid string) (*ScheduleOut, error)
	GetSchedules(ctx context.Context, filters map[string]any) (*SchedulesOut, error)
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

func (s *service) UpdateSchedule(ctx context.Context, uid string, schedule Schedule) error {
	err := s.repo.UpdateSchedule(ctx, uid, schedule)
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
	out.toList(schedules, total)

	return out, nil
}

func (s *service) GetByUID(ctx context.Context, uid string) (*ScheduleOut, error) {
	scheduleEntity, err := s.repo.GetByUID(ctx, uid)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &ScheduleOut{}
	out.PopulateFromEntity(scheduleEntity)

	return out, nil
}

func (s *service) GetSchedules(ctx context.Context, filters map[string]any) (*SchedulesOut, error) {
	schedules, total, err := s.repo.Query(ctx, filters)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &SchedulesOut{}
	out.toList(schedules, total)

	return out, nil
}

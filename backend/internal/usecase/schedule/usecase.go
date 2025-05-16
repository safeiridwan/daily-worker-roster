package user

import (
	"backend/internal/service/middleware"
	"backend/internal/service/schedule"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type Usecase interface {
	CreateSchedule(ctx context.Context, input ScheduleIn, currentUser middleware.Identity) error
}

type usecase struct {
	scheduleService schedule.Service
}

func NewUsecase(scheduleService schedule.Service) Usecase {
	return &usecase{
		scheduleService: scheduleService,
	}
}

var layoutDate = "02-01-2006" // DD-MM-YYYY
var layoutTime = "15:04"

func (u *usecase) CreateSchedule(ctx context.Context, input ScheduleIn, currentUser middleware.Identity) error {
	parsedDate, err := time.Parse(layoutDate, input.Date)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	// Parse time
	startT, err := time.Parse(layoutTime, input.StartTime)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	endT, err := time.Parse(layoutTime, input.EndTime)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	startTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), startT.Hour(), startT.Minute(), 0, 0, time.UTC)
	endTime := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), endT.Hour(), endT.Minute(), 0, 0, time.UTC)

	if !startTime.Before(endTime) {
		err := errors.New("start_time must be before end_time")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	exists, err := u.scheduleService.GetByRangeDate(ctx, parsedDate, startTime, endTime)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	uid, err := uuid.Parse(currentUser.UserUID)
	for _, data := range exists.Data {
		if data.UserUID == uid {
			err := errors.New("schedule already exist")
			zerolog.Ctx(ctx).Error().Err(err).Msg("")
			return err
		}
	}

	if exists.Total > 0 {
		err := errors.New("schedule already assigned")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	scheduleEntity := schedule.Schedule{
		UserUID:   uid,
		Date:      sql.NullTime{Time: parsedDate, Valid: true},
		StartTime: sql.NullTime{Time: startTime, Valid: true},
		EndTime:   sql.NullTime{Time: endTime, Valid: true},
		Status:    sql.NullString{String: input.Status, Valid: true},
		CreatedBy: sql.NullString{String: "System", Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := u.scheduleService.SaveSchedule(ctx, scheduleEntity); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

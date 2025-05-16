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
	EditSchedule(ctx context.Context, input ScheduleIn, uid string, currentUser middleware.Identity) error
	UpdateStatusSchedule(ctx context.Context, uid, status string, currentUser middleware.Identity) error
	GetSchedules(ctx context.Context, filters map[string]any) (*SchedulesOut, error)
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

	if parsedDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		err := errors.New("schedule date cannot be before today")
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
	startLimit := time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC) // 08:00 AM
	endLimit := time.Date(0, 1, 1, 20, 0, 0, 0, time.UTC)  // 08:00 PM

	if !startTime.Before(endTime) {
		err := errors.New("start_time must be before end_time")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if startTime.Before(startLimit) {
		err := errors.New("start time cannot be earlier than 08:00 AM")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if endTime.After(endLimit) {
		err := errors.New("end time cannot be later than 08:00 PM")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

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
		CreatedBy: sql.NullString{String: currentUser.Email, Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := u.scheduleService.SaveSchedule(ctx, scheduleEntity); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (u *usecase) EditSchedule(ctx context.Context, input ScheduleIn, uid string, currentUser middleware.Identity) error {
	exist, err := u.scheduleService.GetByUID(ctx, uid)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if exist == nil {
		err := errors.New("schedule does not exist")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	parsedDate, err := time.Parse(layoutDate, input.Date)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if parsedDate.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		err := errors.New("schedule date cannot be before today")
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
	startLimit := time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC) // 08:00 AM
	endLimit := time.Date(0, 1, 1, 20, 0, 0, 0, time.UTC)  // 08:00 PM

	if !startTime.Before(endTime) {
		err := errors.New("start_time must be before end_time")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if startTime.Before(startLimit) {
		err := errors.New("start time cannot be earlier than 08:00 AM")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if endTime.After(endLimit) {
		err := errors.New("end time cannot be later than 08:00 PM")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	exists, err := u.scheduleService.GetByRangeDate(ctx, parsedDate, startTime, endTime)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	userUid, err := uuid.Parse(currentUser.UserUID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	for _, data := range exists.Data {
		if data.UserUID == userUid {
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

	scheduleEntity := schedule.Schedule{
		UserUID:   userUid,
		Date:      sql.NullTime{Time: parsedDate, Valid: true},
		StartTime: sql.NullTime{Time: startTime, Valid: true},
		EndTime:   sql.NullTime{Time: endTime, Valid: true},
		Status:    sql.NullString{String: input.Status, Valid: true},
		UpdatedBy: sql.NullString{String: currentUser.Email, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := u.scheduleService.UpdateSchedule(ctx, uid, scheduleEntity); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (u *usecase) UpdateStatusSchedule(ctx context.Context, uid, status string, currentUser middleware.Identity) error {
	exist, err := u.scheduleService.GetByUID(ctx, uid)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	if exist == nil {
		err := errors.New("schedule does not exist")
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	scheduleEntity := schedule.Schedule{
		Status:    sql.NullString{String: status, Valid: true},
		UpdatedBy: sql.NullString{String: currentUser.Email, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := u.scheduleService.UpdateSchedule(ctx, uid, scheduleEntity); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (u *usecase) GetSchedules(ctx context.Context, filters map[string]any) (*SchedulesOut, error) {
	schedules, err := u.scheduleService.GetSchedules(ctx, filters)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("")
		return nil, err
	}

	out := &SchedulesOut{}
	out.toList(schedules.Data, schedules.Total)

	return out, nil
}

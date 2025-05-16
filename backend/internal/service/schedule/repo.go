package schedule

import (
	"backend/pkg/gorm/dbcontext"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Schedule struct {
	ID        int64          `gorm:"primaryKey"`
	UID       uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
	UserUID   uuid.UUID      `gorm:"column:user_uid"`
	Date      sql.NullTime   `gorm:"column:date"`
	StartTime sql.NullTime   `gorm:"column:start_time"`
	EndTime   sql.NullTime   `gorm:"column:end_time"`
	Status    sql.NullString `gorm:"column:status"`
	CreatedAt sql.NullTime   `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	CreatedBy sql.NullString `gorm:"column:created_by"`
	UpdatedAt sql.NullTime   `gorm:"column:updated_at"`
	UpdatedBy sql.NullString `gorm:"column:updated_by"`
	DeletedAt sql.NullTime   `gorm:"column:deleted_at"`
	DeletedBy sql.NullString `gorm:"column:deleted_by"`
}

func (Schedule) TableName() string { return "schedule" }

type ScheduleRepository interface {
	SaveSchedule(ctx context.Context, schedule Schedule) error
	GetByRangeDate(ctx context.Context, date, startTime, endTime time.Time) ([]Schedule, int64, error)
}

type scheduleRepository struct {
	db dbcontext.Sessions
}

func NewScheduleRepository(db dbcontext.Sessions) ScheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}

func (r *scheduleRepository) SaveSchedule(ctx context.Context, schedule Schedule) error {
	return r.db.With(ctx).Create(&schedule).Error
}

func (r *scheduleRepository) GetByRangeDate(ctx context.Context, date, startTime, endTime time.Time) ([]Schedule, int64, error) {
	var schedules []Schedule
	err := r.db.With(ctx).
		Model(&Schedule{}).
		Where("date = ? "+
			"AND status = 'APPROVED' "+
			"AND ((start_time < ? AND end_time > ?) "+
			"OR (start_time < ? AND end_time > ?) "+
			"OR (start_time >= ? AND end_time <= ?)) ",
			date, endTime, startTime, startTime, endTime, startTime, endTime).
		Find(&schedules).
		Error

	if err != nil {
		return schedules, 0, err
	}

	return schedules, int64(len(schedules)), nil
}

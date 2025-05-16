package user

import (
	"backend/internal/service/schedule"
	"github.com/google/uuid"
	"time"
)

type ScheduleOut struct {
	ID        int64      `json:"-"`
	UID       uuid.UUID  `json:"uid"`
	UserUID   uuid.UUID  `json:"user_uid"`
	Date      time.Time  `json:"date" `
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy string     `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	DeletedBy string     `json:"deleted_by"`
}

func (o *ScheduleOut) PopulateFromEntity(e *schedule.ScheduleOut) {
	o.ID = e.ID
	o.UID = e.UID
	o.UserUID = e.UserUID
	o.Date = e.Date
	o.StartTime = e.StartTime
	o.EndTime = e.EndTime
	o.CreatedAt = e.CreatedAt
	o.CreatedBy = e.CreatedBy
	if !e.UpdatedAt.IsZero() {
		o.UpdatedAt = e.UpdatedAt
	}
	o.UpdatedBy = e.UpdatedBy
	if !e.DeletedAt.IsZero() {
		o.DeletedAt = e.DeletedAt
	}
	o.DeletedBy = e.DeletedBy
}

type SchedulesOut struct {
	Data  []ScheduleOut
	Total int64
}

func (o *SchedulesOut) toList(m []schedule.ScheduleOut, total int64) {
	o.Data = make([]ScheduleOut, 0, len(m))
	o.Total = total

	for _, schedule := range m {
		var updatedAt *time.Time
		if !schedule.UpdatedAt.IsZero() {
			updatedAt = schedule.UpdatedAt
		}

		var deletedAt *time.Time
		if !schedule.DeletedAt.IsZero() {
			deletedAt = schedule.DeletedAt
		}

		o.Data = append(o.Data, ScheduleOut{
			UID:       schedule.UID,
			UserUID:   schedule.UserUID,
			Date:      schedule.Date,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
			CreatedAt: schedule.CreatedAt,
			CreatedBy: schedule.CreatedBy,
			UpdatedAt: updatedAt,
			UpdatedBy: schedule.UpdatedBy,
			DeletedAt: deletedAt,
			DeletedBy: schedule.DeletedBy,
		})
	}
}

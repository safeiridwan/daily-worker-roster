package schedule

import (
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

type SchedulesOut struct {
	Data  []ScheduleOut
	Total int64
}

func (o *SchedulesOut) PopulateFromEntity(m []Schedule, total int64) {
	o.Data = make([]ScheduleOut, 0, len(m))
	o.Total = total

	for _, schedule := range m {
		var updatedAt *time.Time
		if !schedule.UpdatedAt.Time.IsZero() {
			updatedAt = &schedule.UpdatedAt.Time
		}

		var deletedAt *time.Time
		if !schedule.DeletedAt.Time.IsZero() {
			deletedAt = &schedule.DeletedAt.Time
		}

		o.Data = append(o.Data, ScheduleOut{
			UID:       schedule.UID,
			UserUID:   schedule.UserUID,
			Date:      schedule.Date.Time,
			StartTime: schedule.StartTime.Time,
			EndTime:   schedule.EndTime.Time,
			CreatedAt: schedule.CreatedAt.Time,
			CreatedBy: schedule.CreatedBy.String,
			UpdatedAt: updatedAt,
			UpdatedBy: schedule.UpdatedBy.String,
			DeletedAt: deletedAt,
			DeletedBy: schedule.DeletedBy.String,
		})
	}
}

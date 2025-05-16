package user

import (
	"backend/pkg/chi/request"
	"net/http"
)

type ScheduleIn struct {
	Date      string `json:"date"`
	StartTime string `json:"start_time" `
	EndTime   string `json:"end_time" `
	Status    string `json:"status"`
}

func (c ScheduleIn) Bind(_ *http.Request) error {
	return request.ValidateStruct(c)
}

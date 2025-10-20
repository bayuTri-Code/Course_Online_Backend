package dto

import "time"

type ActivityResponse struct {
	ID           string        `json:"id"`
	ActivityName string        `json:"activity_name"`
	When         time.Time     `json:"when"`
	User         UserSimpleDTO `json:"user"`
}

type UserSimpleDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	Biodata  *BiodataResponseForAct `json:"biodata,omitempty"`
}

type BaseResponseActivityByUserId struct {
	Status  bool            `json:"status" example:"true"`
	Message string          `json:"message"`
	Data    ActivityResponse `json:"data"`
}

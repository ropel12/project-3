package event

import (
	entity "github.com/ropel12/project-3/app/entities"
)

type (
	ResponseEvent struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		Date         string `json:"date"`
		Location     string `json:"location"`
		EndDate      string `json:"end_date"`
		HostedBy     string `json:"hosted_by"`
		Image        string `json:"image"`
		Participants int    `json:"participants"`
	}
	ReqCreate struct {
		Name      string  `form:"name" validate:"required"`
		StartDate string  `form:"date" validate:"required"`
		Duration  float32 `form:"duration" validate:"required"`
		Details   string  `form:"details" validate:"required"`
		Quota     int     `form:"quota" validate:"required"`
		HostedBy  string  `form:"hosted_by" validate:"required"`
		Location  string  `form:"location" validate:"required"`
		Rtype     string  `form:"type" json:"type" validate:"required"`
		Types     []entity.Type
		Image     string
		Uid       int
	}
	Response struct {
		Limit     int `json:"limit"`
		Page      int `json:"page"`
		TotalPage int `json:"total_page"`
		TotalData int `json:"total_data"`
		Data      any `json:"data"`
	}

	UserComments struct {
		Name    string
		Image   string
		Comment string
	}
	UserParticipant struct {
		Name  string
		Image string
	}
	DetailEvent struct {
		Id           int               `json:"id"`
		Name         string            `json:"name"`
		Date         string            `json:"date"`
		Details      string            `json:"details"`
		Location     string            `json:"location"`
		Duration     float32           `json:"duration"`
		HostedBy     string            `json:"hosted_by"`
		Quota        int               `json:"quota"`
		Image        string            `json:"image"`
		Types        []entity.Type     `json:"types"`
		Participants []UserParticipant `json:"participants"`
		UserComments []UserComments    `json:"comments"`
	}
	ResponseDetailEvent struct {
		Data any `json:"data"`
	}
)

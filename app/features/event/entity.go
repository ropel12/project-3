package event

import (
	entity "github.com/ropel12/project-3/app/entities"
)

type (
	ResponseEvent struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		Date         string `json:"date"`
		Quota        int    `json:"quota"`
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
	ReqUpdate struct {
		Id        uint    `form:"id" validate:"required"`
		Name      string  `form:"name" validate:"required"`
		StartDate string  `form:"date" validate:"required"`
		Duration  float32 `form:"duration" validate:"required"`
		Details   string  `form:"details" validate:"required"`
		Quota     int     `form:"quota" validate:"required"`
		HostedBy  string  `form:"hosted_by" validate:"required"`
		Location  string  `form:"location" validate:"required"`
		Rtype     string  `form:"type" validate:"required"`
		Image     string  `form:"image"`
		Types     []TypeEvent
	}
	TypeEvent struct {
		Id       int    `json:"id"`
		TypeName string `json:"type_name"`
		Price    int    `json:"price"`
	}
	ReqCreateComment struct {
		EventId int    `json:"event_id" validate:"required"`
		Comment string `json:"comment" validate:"required"`
		Uid     int
	}
	Response struct {
		Limit     int `json:"limit,omitempty"`
		Page      int `json:"page,omitempty"`
		TotalPage int `json:"total_page,omitempty"`
		TotalData int `json:"total_data,omitempty"`
		Data      any `json:"data"`
	}

	UserComments struct {
		Name    string `json:"name,omitempty"`
		Image   string `json:"image,omitempty"`
		Comment string `json:"comment,omitempty"`
	}
	UserParticipant struct {
		Name  string `json:"name,omitempty"`
		Image string `json:"image,omitempty"`
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
		Types        []TypeEvent       `json:"types"`
		Participants []UserParticipant `json:"participants"`
		UserComments []UserComments    `json:"comments"`
	}
	ResponseDetailEvent struct {
		Data any `json:"data"`
	}
)

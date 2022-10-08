package models

import "github.com/gofrs/uuid"

// API response model for SchedulesAPI
type SchedulesAPIResponse struct {
	Page       uint       `json:"page"`
	PageSize   uint       `json:"page_size"`
	TotalCount uint       `json:"total_count"`
	Schedules  []Schedule `json:"schedules,omitempty"`
}

// Schedule is the struct for representing schedules
type Schedule struct {
	ID uuid.UUID `json:"id,omitempty"`

	// CronExpression is the UNIX cron expression (quartz expression)
	//
	// Example:
	// 	0 15 5 ? * WED,SUN *
	CronExpression string `json:"cron_expression,omitempty"`
}

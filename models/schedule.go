package models

import "github.com/gofrs/uuid"

// Schedule is the struct for representing schedules
type Schedule struct {
	ID uuid.UUID `json:"id,omitempty"`

	// CronExpression is the UNIX cron expression (quartz expression)
	//
	// Example:
	// 	0 15 5 ? * WED,SUN *
	CronExpression string `json:"cron_expression,omitempty"`
}

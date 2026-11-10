package models

import "github.com/meshery/schemas/models/v1beta2/schedule"

type Schedule = schedule.Schedule

// SchedulesAPIResponse is the API response for schedules list endpoint.
type SchedulesAPIResponse = schedule.SchedulePage

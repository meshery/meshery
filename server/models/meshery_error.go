package models

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

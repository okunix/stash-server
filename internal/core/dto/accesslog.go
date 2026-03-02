package dto

import "github.com/google/uuid"

type AccessLogResponse struct {
}

type ListAccessLogRequest struct {
	StashID uuid.UUID `json:"stash_id"`
	Limit   uint      `json:"limit"`
	Offset  uint      `json:"offset"`
}

type AddLogEntryRequest struct{}

package dto

import (
	"time"

	"github.com/google/uuid"
)

type SecretResponse struct {
	Data       map[string]string `json:"data"`
	UnlockedAt time.Time         `json:"unlocked_at"`
}

type GetSecretByStashID struct {
	UserID  uuid.UUID `json:"user_id"`
	StashID uuid.UUID `json:"stash_id"`
}

type ListSecretResponse struct {
	Page    *Page            `json:"page,omitempty"`
	Content []SecretResponse `json:"content"`
}

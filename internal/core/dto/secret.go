package dto

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/stash-password-manager/stash-server/internal/core/domain/secret"
)

type SecretResponse struct {
	//Data       map[string]string `json:"data"`
	Keys       []string  `json:"keys"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

type GetSecretByStashID struct {
	UserID  uuid.UUID `json:"user_id"`
	StashID uuid.UUID `json:"stash_id"`
}

type ListSecretResponse struct {
	Page    *Page            `json:"page,omitempty"`
	Content []SecretResponse `json:"content"`
}

type AddSecret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (a AddSecret) Validate() (map[string]string, bool) {
	problems := make(map[string]string)
	if err := secret.ValidateEntryName(a.Name); err != nil {
		problems["name"] = err.Error()
	}
	return problems, len(problems) == 0
}

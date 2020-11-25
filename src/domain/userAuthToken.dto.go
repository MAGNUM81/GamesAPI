package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"strings"
)

type UserSession struct {
	Token string `json:"token"`
	UserId uint64 `json:"user_id"`
	ExpiresAt int64 `json:"expires_at"`
}

func (t *UserSession) Validate() errorUtils.EntityError {
	if strings.Trim(t.Token, " ") == "" {
		return errorUtils.NewUnprocessableEntityError("Token cannot be empty")
	}
	return nil
}
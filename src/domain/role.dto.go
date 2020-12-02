package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"time"
)

type UserRole struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	UserID    uint64     `gorm:"column:user_id" json:"user_id"`
	Name      string     `gorm:"column:roleName" json:"name"`
}

func (r *UserRole) Validate() errorUtils.EntityError {
	if r.UserID <= 0 {
		return errorUtils.NewUnprocessableEntityError("Role UserID is invalid")
	}

	if r.Name == "" {
		return errorUtils.NewUnprocessableEntityError("Role Name cannot be empty")
	}
	return nil
}

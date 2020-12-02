package domain

import (
	"GamesAPI/src/utils/errorUtils"
	"regexp"
	"time"
)

type User struct {
	ID           uint64     `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
	Name         string     `gorm:"column:name;not null;" json:"name"`
	Email        string     `gorm:"column:email;not null;unique" json:"email"`
	PasswordHash string     `gorm:"column:password_hash;not null;default:'hashpass'" json:"password_hash"`
	SteamUserId  string     `gorm:"column:steam_user_id;not null" json:"steam_user_id"`
}

func (u *User) Validate() errorUtils.EntityError {

	//check for empty name
	if u.Name == "" {
		return errorUtils.NewUnprocessableEntityError("User name cannot be empty")
	}

	//check for empty email
	if u.Email == "" {
		return errorUtils.NewUnprocessableEntityError("User email cannot be empty")
	}

	//check for bad email format.
	//inspired from https://www.golangprograms.com/regular-expression-to-validate-email-address.html
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(u.Email) {
		return errorUtils.NewUnprocessableEntityError("User email is not formatted correctly")
	}

	return nil
}

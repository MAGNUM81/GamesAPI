package domain

import "github.com/jinzhu/gorm"

func InitRepositories(db *gorm.DB) {
	UserRepo.Initialize(db)
	GameRepo.Initialize(db)
}

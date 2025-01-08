package model

import (
	"gorm.io/gorm"
)

func ModelMigration(db *gorm.DB){
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Book{})
	// db.AutoMigrate(&BorrowedBook{})
}
package model

import (
	"gorm.io/gorm"
)

func ModelMigration(db *gorm.DB){
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Book{})
	db.Exec("ALTER TABLE books ADD CONSTRAINT unique_user_isbn UNIQUE (isbn, user_id);")
	// db.AutoMigrate(&BorrowedBook{})
}
package model

import "time"

type User struct {
	ID                int       `gorm:"primaryKey;autoIncrement"`
	Email             string    `gorm:"size:255;unique;not null"`
	Username          string    `gorm:"size:255;not null"`
	Password          string    `gorm:"size:255;not null"`
	Role              string    `gorm:"type:varchar(20);check(role in ('user','admin'));default:'user';not null"`
	Status            string    `gorm:"type:varchar(20);check(status in ('active','inactive'));default:'active';not null"`
	IsEmailVerified   bool      `gorm:"default:false"`
	VerificationToken string    `gorm:"size:255"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

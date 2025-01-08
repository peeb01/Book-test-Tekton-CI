package model

import "time"

type Book struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    ISBN        string    `gorm:"size:13;unique;not null"`
    Title       string    `gorm:"size:255;not null"`
    Author      string    `gorm:"size:255;not null"`
    Price       float64   `gorm:"type:decimal(10,2);not null"`
    Stock       int       `gorm:"not null"`
    PublishedAt time.Time `gorm:"type:date"`
	ImageLink	string
	BookType	string	  `gorm:"size:50;not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

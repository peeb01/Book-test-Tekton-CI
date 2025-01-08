package model

import "time"

type Book struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    ISBN        string    `gorm:"size:13;unique;not null" json:"ISBN"`
    Title       string    `gorm:"size:255;not null" json:"Title"`
    Author      string    `gorm:"size:255;not null" json:"Author"`
    Price       float64   `gorm:"type:decimal(10,2);not null" json:"Price"`
    Stock       int       `gorm:"not null" json:"Stock"`
    PublishedAt time.Time `gorm:"type:date" json:"PublishedAt"`
    ImageLink   string    `json:"ImageLink"`
    BookType    string    `gorm:"size:50;not null" json:"BookType"`
    CreatedAt   time.Time `json:"-"`
    UpdatedAt   time.Time `json:"-"`
}

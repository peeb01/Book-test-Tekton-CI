package model

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ISBN        string    `gorm:"size:13;not null" json:"ISBN"`
	Title       string    `gorm:"size:255;not null" json:"Title"`
	Author      string    `gorm:"size:255;not null" json:"Author"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"Price"`
	Stock       int       `gorm:"not null" json:"Stock"`
	PublishedAt time.Time `gorm:"type:date" json:"PublishedAt"`
	ImageLink   string    `json:"ImageLink"`
	BookType    string    `gorm:"size:50;not null" json:"BookType"`
	UserID      uint      `gorm:"not null" json:"UserID"` // Foreign Key to User
	Users       []User    `gorm:"many2many:user_books;"`  // Many-to-Many relationship
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

// Add unique composite constraint for ISBN + UserID
func (Book) TableName() string {
	return "books"
}
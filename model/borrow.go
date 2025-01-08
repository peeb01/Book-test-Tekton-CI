package model

import "time"

type BorrowedBook struct {
    ID          uint      `gorm:"primaryKey;autoIncrement"`
    UserID      uint      `gorm:"not null"` // Foreign Key จาก user
    BookID      uint      `gorm:"not null"` // Foreign Key จาก book
    BorrowDate  time.Time `gorm:"type:date;not null"`
    ReturnDate  time.Time `gorm:"type:date;not null"`
    Status      string    `gorm:"type:varchar(20);check(status in ('borrowed','returned','late'));default:'borrowed';not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time

    // ความสัมพันธ์
    User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Book        Book      `gorm:"foreignKey:BookID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
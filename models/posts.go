package models

import "time"

type Posts struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"size:200;not null"`
	Content     string    `gorm:"type:text;not null"`
	Category    string    `gorm:"size:100;not null"`
	CreatedDate time.Time `gorm:"autoCreateTime"`
	UpdatedDate time.Time `gorm:"autoUpdateTime"`
	Status      string    `gorm:"size:100;not null;check:status IN ('Publish', 'Draft', 'Thrash')"`
}

type Users struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"size:100;not null;unique"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:100;not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

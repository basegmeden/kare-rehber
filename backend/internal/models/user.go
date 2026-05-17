package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin       Role = "admin"
	RoleCoach       Role = "koc"
	RoleCoordinator Role = "koordinator"
	RoleParent      Role = "veli"
	RoleStudent     Role = "ogrenci"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "active"
	UserStatusPassive UserStatus = "passive"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `gorm:"not null" json:"name"`
	Surname      string         `gorm:"not null" json:"surname"`
	Phone        string         `gorm:"uniqueIndex;not null" json:"phone"`
	DateOfBirth  *time.Time     `json:"date_of_birth,omitempty"`
	City         string         `json:"city"`
	Role         Role           `gorm:"not null" json:"role"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Status       UserStatus     `gorm:"default:active" json:"status"`
}

func (u *User) FullName() string {
	return u.Name + " " + u.Surname
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID             string         `json:"id" gorm:"type:char(19);primaryKey"`
	OrganizationID string         `json:"-" gorm:"type:char(19);index;not null"`
	Title          string         `json:"title" gorm:"type:varchar(64)"`
	CreatedAt      time.Time      `json:"createdAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type Permission struct {
	ID        string    `json:"id" gorm:"type:char(19);primaryKey"`
	RoleID    string    `json:"roleId" gorm:"type:char(19);index;not null"`
	Resource  string    `json:"resource" gorm:"type:varchar(128)"`
	Action    string    `json:"action" gorm:"type:varchar(64)"`
	CreatedAt time.Time `json:"createdAt"`
}

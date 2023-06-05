package model

import "time"

type Task struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Done        bool      `json:"done" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Deadline    string    `json:"deadline,omitempty" gorm:"timestamp;not null;default: (now() + interval '2 day')"`
	UserId      int       `json:"-"`
	Username    string    `json:"username" gorm:"-"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"autoUpdateTime"`
	DeletedAt   time.Time `json:"-"`
	User        User      `json:"-" gorm:"foreignKey:UserId"`
}

type User struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Password  string    `json:"password" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"-"`
}

type Tasks []Task

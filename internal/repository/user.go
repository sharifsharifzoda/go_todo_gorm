package repository

import (
	"gorm.io/gorm"
	"todo_gorm/model"
)

type AuthPostgres struct {
	db *gorm.DB
}

func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *model.User) (int, error) {
	tx := r.db.Create(&user)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return user.Id, nil
}

func (r *AuthPostgres) GetUser(email string) (user model.User, err error) {
	tx := r.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return model.User{}, tx.Error
	}

	return user, nil
}

func (r *AuthPostgres) IsEmailUsed(email string) bool {
	var user model.User
	tx := r.db.Where("email = ?", email).Find(&user)
	if tx.Error != nil {
		return false
	}

	if user.Email == "" {
		return false
	}

	return true
}

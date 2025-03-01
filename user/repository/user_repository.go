package repository

import (
	"user_service/types"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (u *UserRepository) Register(payload types.RegisterPayload) error {
	if err := u.DB.Create(&types.User{
		Email:    payload.Email,
		Password: payload.Password,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (u *UserRepository) GetUserByID(ID int) (types.User, error) {
	var user types.User
	if err := u.DB.Where("id = ?", ID).First(&user).Error; err != nil {
		return types.User{}, err
	}
	return user, nil
}

package repository

import (
	"user_service/types"
	"user_service/user/usergrpc"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(*usergrpc.RegisterPayload) error
	GetUserByEmail(string) (*usergrpc.User, error)
	GetUserByID(int) (*usergrpc.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (u *UserRepositoryImpl) Register(payload *usergrpc.RegisterPayload) error {
	if err := u.DB.Create(&usergrpc.User{
		Email:    payload.Email,
		Password: payload.Password,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) GetUserByEmail(email string) (*usergrpc.User, error) {
	var user types.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &usergrpc.User{
		ID:        int32(user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *UserRepositoryImpl) GetUserByID(ID int) (*usergrpc.User, error) {
	var user types.User
	if err := u.DB.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &usergrpc.User{
		ID:        int32(user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

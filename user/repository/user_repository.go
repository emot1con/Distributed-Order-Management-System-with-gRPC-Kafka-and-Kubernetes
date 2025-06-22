package repository

import (
	"user_service/proto"
	"user_service/types"

	"gorm.io/gorm"
)

type UserRepository interface {
	Register(*proto.RegisterPayload) error
	GetUserByEmail(string) (*proto.User, error)
	GetUserByID(int) (*proto.User, error)
	UpdateUser(user *proto.User) error
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (u *UserRepositoryImpl) Register(payload *proto.RegisterPayload) error {
	if err := u.DB.Create(&proto.User{
		Email:      payload.Email,
		Password:   payload.Password,
		FullName:   payload.FullName,
		Provider:   payload.Provider,
		ProviderId: payload.ProviderId,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) GetUserByEmail(email string) (*proto.User, error) {
	var user types.User
	if err := u.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &proto.User{
		ID:        int32(user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *UserRepositoryImpl) GetUserByID(ID int) (*proto.User, error) {
	var user types.User
	if err := u.DB.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &proto.User{
		ID:        int32(user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *UserRepositoryImpl) UpdateUser(user *proto.User) error {
	// Update the user in database
	if err := u.DB.Model(&types.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"full_name": user.FullName,
		"email":     user.Email,
		"password":  user.Password,
	}).Error; err != nil {
		return err
	}

	return nil
}

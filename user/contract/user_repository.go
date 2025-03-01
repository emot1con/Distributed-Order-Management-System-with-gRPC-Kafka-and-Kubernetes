package contract

import "user_service/types"

type UserRepository interface {
	Register(types.RegisterPayload) error
	GetUserByEmail(string) (types.User, error)
	GetUserByID(int) (types.User, error)
}

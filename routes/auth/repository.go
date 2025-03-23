package auth

import "github.com/yantology/retail-pro-be/pkg/customerror"

type AuthRepository struct {
	db AuthDBInterface
}

func NewAuthRepository(db AuthDBInterface) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) CheckExistingEmail(email string) *customerror.CustomError {
	return ar.db.CheckExistingEmail(email)
}

func (ar *AuthRepository) SaveActivationToken(req *ActivationTokenRequest) *customerror.CustomError {
	return ar.db.SaveActivationToken(req)
}

func (ar *AuthRepository) GetActivationToken(req *GetActivationTokenRequest) (string, *customerror.CustomError) {
	return ar.db.GetActivationToken(req)
}

func (ar *AuthRepository) CreateUser(req *CreateUserRequest) *customerror.CustomError {
	return ar.db.CreateUser(req)
}

func (ar *AuthRepository) GetUserByEmail(email string) (*User, *customerror.CustomError) {
	return ar.db.GetUserByEmail(email)
}

func (ar *AuthRepository) UpdateUserPassword(req *UpdatePasswordRequest) *customerror.CustomError {
	return ar.db.UpdateUserPassword(req)
}

package service

import (
	"golang.org/x/crypto/bcrypt"
	"gtihub.com/raditsoic/telkom-storage-ms/src/database/repository"
	"gtihub.com/raditsoic/telkom-storage-ms/src/model"
	"gtihub.com/raditsoic/telkom-storage-ms/src/utils"
)

type AuthService struct {
	repository repository.AdminRepository
	jwtUtils   *utils.JWTUtils
}

func NewAuthService(repo repository.AdminRepository, jwtUtils *utils.JWTUtils) *AuthService {
	return &AuthService{
		repository: repo,
		jwtUtils:   jwtUtils,
	}
}

func (service *AuthService) AdminLogin(LoginReq *model.LoginRequest) (*model.LoginResponse, error) {
	admin, err := service.repository.GetAdminByUsername(LoginReq.Username)
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(LoginReq.Password))
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	token, err := service.jwtUtils.GenerateJWT(admin.ID)
	if err != nil {
		return nil, err
	}

	response := &model.LoginResponse{
		Token:    token,
		Username: admin.Username,
		Message:  "Login successful",
	}

	return response, nil
}

func (service *AuthService) AdminRegister(RegisReq *model.RegisterRequest) (*model.RegisterResponse, error) {
	_, err := service.repository.GetAdminByUsername(RegisReq.Username)
	if err == nil {
		return nil, utils.ErrUsernameExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RegisReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{
		Username: RegisReq.Username,
		Password: string(hashedPassword),
	}

	err = service.repository.RegisterAdmin(admin)
	if err != nil {
		return nil, err
	}

	response := &model.RegisterResponse{
		Username: admin.Username,
		Message:  "Admin created",
	}

	return response, nil
}

func (service *AuthService) GetAdmin() ([]model.Admin, error) {
	admins, err := service.repository.GetAdmins()
	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (service *AuthService) DeleteAdmin(id string) (*model.DeleteResponse, error) {
	if id == "" {
		return nil, utils.ErrInvalidID
	}

	err := service.repository.DeleteAdmin(id)
	if err != nil {
		return nil, err
	}

	response := &model.DeleteResponse{
		Message: "Admin deleted",
	}

	return response, nil
}

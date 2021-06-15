package services

import (
	vo "chatapp/internal/app/valobjects"
	entity "chatapp/internal/data/entities"
	repository "chatapp/internal/data/repositories"
	db "chatapp/internal/database"
)

type RegistrationService struct {
	dto *RegistrationDetailsDTO
}

// New User's Registration Details DTO
type RegistrationDetailsDTO struct {
	Name     vo.Name
	Username vo.Email
	Password string
}

// CReate a new instance
func NewRegistrationService(dto *RegistrationDetailsDTO) ApplicationServiceInterface {
	this := &RegistrationService{dto: dto}
	return this
}

// Executes the service
func (this *RegistrationService) Execute() error {
	userRepo := repository.NewUserRepository(db.GetInstance())
	newUser := entity.User{
		Name:     this.dto.Name,
		Username: this.dto.Username,
		Password: this.dto.Password,
	}
	userRepo.Add(&newUser)
	return nil
}

package services

import (
	repository "chatapp/internal/data/repositories"
	db "chatapp/internal/database"
	"chatapp/internal/lib/util"

	"github.com/markbates/goth"
)

type AuthenticationService struct {
	dto      *CredentialsDTO
	gothUser *goth.User
}

// Requested User's credentials DTO
type CredentialsDTO struct {
	Username string
	Password string
}

// CReate a new instance
func NewAuthenticationService(dto *CredentialsDTO, gothUser *goth.User) ApplicationServiceInterface {
	this := &AuthenticationService{
		dto:      dto,
		gothUser: gothUser,
	}
	return this
}

// LoginHandler handles application's login process
func (this *AuthenticationService) Execute() error {
	userRepo := repository.NewUserRepository(db.GetInstance())
	targetedUser := userRepo.FindByUsername(this.dto.Username)
	// If a password exists for the given user AND if it is the same as the password
	// we received, we can move ahead, otherwise return an "Unauthorized" status
	ok, err := util.ComparePasswords(this.dto.Password, targetedUser.Password)
	if !ok || err != nil {
		return err
	}
	// Create JWT token
	token, err := util.CreateJWTAccessToken(this.dto.Username)
	if err != nil {
		return err
	}
	// Update the goth user for caller
	this.gothUser.Name = targetedUser.Name.String()
	this.gothUser.Email = targetedUser.Email().String()
	this.gothUser.AccessToken = token
	return nil
}

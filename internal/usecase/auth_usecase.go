package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/l0ng7h0r/ecommerce/internal/domain"
	"github.com/l0ng7h0r/ecommerce/internal/repository"
	"github.com/l0ng7h0r/ecommerce/pkg/config"
	"github.com/l0ng7h0r/ecommerce/pkg/security"
)

type AuthUsecase struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo *repository.UserRepository, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *AuthUsecase) Register(req *domain.RegisterReq) (*domain.AuthResponse, error) {
	existing, _ := u.userRepo.GetUserByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.CreateUser(req.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	roleName := req.Role
	if roleName == "" {
		roleName = "user"
	}
	if err := u.userRepo.AssignRole(user.ID, roleName); err != nil {
		return nil, err
	}

	roles, err := u.userRepo.GetUserRoles(user.ID)
	if err != nil {
		roles = []string{roleName}
	}
	user.Roles = roles

	accessToken, err := security.GenerateToken(user.ID, user.Email, user.Roles, u.cfg.JWTSecret, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	if err := u.userRepo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		Roles:        user.Roles,
	}, nil
}

func (u *AuthUsecase) Login(req *domain.LoginReq) (*domain.AuthResponse, error) {
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !security.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := security.GenerateToken(user.ID, user.Email, user.Roles, u.cfg.JWTSecret, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	if err := u.userRepo.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		Roles:        user.Roles,
	}, nil
}

func (u *AuthUsecase) Refresh(refreshToken string) (*domain.AuthResponse, error) {
	userID, err := u.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Rotate refresh token
	_ = u.userRepo.DeleteRefreshToken(refreshToken)

	newAccessToken, err := security.GenerateToken(user.ID, user.Email, user.Roles, u.cfg.JWTSecret, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	newRefreshToken := uuid.New().String()
	if err := u.userRepo.SaveRefreshToken(user.ID, newRefreshToken); err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		Roles:        user.Roles,
	}, nil
}

func (u *AuthUsecase) Logout(refreshToken string) error {
	return u.userRepo.DeleteRefreshToken(refreshToken)
}

func (u *AuthUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *AuthUsecase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.GetUserByID(id)
}

func (u *AuthUsecase) DeleteUser(id string) error {
	return u.userRepo.DeleteUser(id)
}

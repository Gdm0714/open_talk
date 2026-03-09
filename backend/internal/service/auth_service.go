package service

import (
	"errors"
	"time"

	"github.com/godongmin/open_talk/backend/internal/config"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/godongmin/open_talk/backend/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidNickname    = errors.New("nickname must be 2-20 characters")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrWrongPassword      = errors.New("current password is incorrect")
)

type AuthService interface {
	Register(email, password, nickname string) (*model.User, string, error)
	Login(email, password string) (*model.User, string, error)
	RefreshToken(tokenString string) (string, error)
	ChangePassword(userID, oldPassword, newPassword string) error
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(email, password, nickname string) (*model.User, string, error) {
	if !validator.ValidateEmail(email) {
		return nil, "", ErrInvalidEmail
	}
	if !validator.ValidatePassword(password) {
		return nil, "", ErrInvalidPassword
	}
	if !validator.ValidateNickname(nickname) {
		return nil, "", ErrInvalidNickname
	}

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, "", ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (*model.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) RefreshToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}

	newToken, err := s.generateToken(claims.Subject)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *authService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrWrongPassword
	}

	if !validator.ValidatePassword(newPassword) {
		return ErrInvalidPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	return s.userRepo.Update(user)
}

func (s *authService) generateToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

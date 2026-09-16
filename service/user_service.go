package service

import (
	"pulse/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(email, password, userName string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return s.repo.RegisterUser(userName, email, string(hashedPassword))
}

func (s *UserService) LoginUser(email, password string) (int64, string, error) {
	user, err := s.repo.LoginUser(email)
	if err != nil {
		return 0, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return 0, "", err
	}
	return user.ID, user.UserName, nil
}

func (s *UserService) Me(id int64) (repository.Me, error) {
	return s.repo.Me(id)
}

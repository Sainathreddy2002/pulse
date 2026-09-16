package service

import (
	"errors"
	"pulse/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrCantFollowYourself       = errors.New("can't follow yourself")
	ErrUserOrFollowerNotFound   = errors.New("user or follower doesn't exist")
	ErrAlreadyFollowing         = errors.New("already following")
)

type FollowService struct {
	follows *repository.FollowRepository
	users   *repository.UserRepository
}

func NewFollowService(follows *repository.FollowRepository, users *repository.UserRepository) *FollowService {
	return &FollowService{follows: follows, users: users}
}

func (s *FollowService) FollowUser(followingID, followerID int64) error {
	if followingID == followerID {
		return ErrCantFollowYourself
	}

	exists, err := s.users.UserExists(followingID)
	if err != nil {
		return err
	}
	followerExists, err := s.users.UserExists(followerID)
	if err != nil {
		return err
	}
	if !exists || !followerExists {
		return ErrUserOrFollowerNotFound
	}

	err = s.follows.FollowUser(followingID, followerID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyFollowing
		}
		return err
	}
	return nil
}

func (s *FollowService) UnfollowUser(followingID, followerID int64) error {
	exists, err := s.users.UserExists(followingID)
	if err != nil {
		return err
	}
	followerExists, err := s.users.UserExists(followerID)
	if err != nil {
		return err
	}
	if !exists || !followerExists {
		return ErrUserOrFollowerNotFound
	}
	return s.follows.UnfollowUser(followingID, followerID)
}

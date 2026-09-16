package service

import (
	"database/sql"
	"errors"
	"pulse/repository"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrCantFollowYourself     = errors.New("can't follow yourself")
	ErrUserOrFollowerNotFound = errors.New("user or follower doesn't exist")
	ErrAlreadyFollowing       = errors.New("already following")
)

type FollowService struct {
	follows       *repository.FollowRepository
	users         *repository.UserRepository
	notifications *repository.NotificationRepository
	outbox        *repository.OutboxRepository
	db            *sql.DB
}

func NewFollowService(follows *repository.FollowRepository, users *repository.UserRepository, notifications *repository.NotificationRepository, outbox *repository.OutboxRepository, db *sql.DB) *FollowService {
	return &FollowService{follows: follows, users: users, notifications: notifications, outbox: outbox, db: db}
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
	tx, txErr := s.db.Begin()
	if txErr != nil {
		return txErr
	}
	err = s.follows.FollowUser(followingID, followerID, tx)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			tx.Rollback()
			return ErrAlreadyFollowing
		}
		tx.Rollback()
		return err
	}
	notifyErr := s.notifications.AddFollowNotification(tx, uuid.New().String(), followerID, followingID, "follow")
	if notifyErr != nil {
		tx.Rollback()
		return notifyErr
	}
	outboxErr := s.outbox.AddProcess(tx, "follow", followerID, followingID)
	if outboxErr != nil {
		tx.Rollback()
		return outboxErr
	}
	return tx.Commit()

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

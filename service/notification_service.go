package service

import "pulse/repository"

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUserNotifications(userID int64) ([]repository.SingleNotification, error) {
	return s.repo.GetUserNotifications(userID)
}

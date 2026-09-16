package repository

import "database/sql"

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (repo *NotificationRepository) AddFollowNotification(q DBTX, id string, caused_by, received_by int64, notifType string) error {
	_, err := q.Exec(`insert into notifications (id,type,caused_by,received_by) values ($1,$2,$3,$4)`, id, notifType, caused_by, received_by)
	return err
}

// func (repo *NotificationRepository) DeleteFollowNotification(q DBTX, caused_by, received_by int64, desc string) error {
// 	_, err := q.Exec(`DELETE FROM notifications WHERE type=$1 AND caused_by=$2 AND received_by=$3)`, desc, caused_by, received_by)
// 	return err
// }

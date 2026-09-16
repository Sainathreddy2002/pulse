package repository

import (
	"database/sql"
	"time"
)

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

type SingleNotification struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	CausedBy   int64     `json:"causedBy"`
	ReceivedBy int64     `json:"receivedBy"`
	CreatedAt  time.Time `json:"createdAt"`
	IsRead     bool      `json:"isRead"`
}

func (repo *NotificationRepository) GetUserNotifications(userID int64) ([]SingleNotification, error) {
	rows, err := repo.db.Query(
		`SELECT id, type, caused_by, received_by, created_at, is_read
		 FROM notifications
		 WHERE received_by = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return []SingleNotification{}, err
	}
	defer rows.Close()

	res := []SingleNotification{}
	for rows.Next() {
		var ele SingleNotification
		if scanErr := rows.Scan(&ele.ID, &ele.Type, &ele.CausedBy, &ele.ReceivedBy, &ele.CreatedAt, &ele.IsRead); scanErr != nil {
			return []SingleNotification{}, scanErr
		}
		res = append(res, ele)
	}
	if err := rows.Err(); err != nil {
		return []SingleNotification{}, err
	}
	return res, nil
}

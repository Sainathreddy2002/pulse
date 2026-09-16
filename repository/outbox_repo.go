package repository

import (
	"database/sql"
	"time"
)

type OutboxRepository struct {
	db *sql.DB
}

func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (repo *OutboxRepository) AddProcess(q DBTX, processType string, causedBy, receivedBy int64) error {
	_, err := q.Exec(`insert into outbox (type,caused_by,received_by) values ($1,$2,$3)`, processType, causedBy, receivedBy)
	return err
}

type SingleProcess struct {
	ID         int64     `json:"id"`
	Type       string    `json:"type"`
	CausedBy   int64     `json:"causedBy"`
	ReceivedBy int64     `json:"receivedBy"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"`
}

func (repo *OutboxRepository) GetProcesses(limit int) ([]SingleProcess, error) {
	tx, txErr := repo.db.Begin()
	if txErr != nil {
		return []SingleProcess{}, txErr
	}
	rows, err := tx.Query(`UPDATE outbox SET status='processing'
WHERE id IN (
  SELECT id FROM outbox
  WHERE processed_at IS NULL and status='pending'
  ORDER BY id
  LIMIT $1
  FOR UPDATE SKIP LOCKED
)
RETURNING id,type,caused_by,received_by,created_at,status`, limit)
	if err != nil {
		tx.Rollback()
		return []SingleProcess{}, err
	}
	var res []SingleProcess
	defer rows.Close()
	for rows.Next() {
		var ele SingleProcess
		if scanErr := rows.Scan(&ele.ID, &ele.Type, &ele.CausedBy, &ele.ReceivedBy, &ele.CreatedAt, &ele.Status); scanErr != nil {
			tx.Rollback()
			return []SingleProcess{}, scanErr
		}
		res = append(res, ele)
	}
	if err := rows.Err(); err != nil {
		tx.Rollback()
		return []SingleProcess{}, err
	}
	rows.Close()
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *OutboxRepository) MarkProcessed(id int64) error {
	_, err := repo.db.Exec(`update outbox set processed_at=NOW(),status='done' where id=$1`, id)
	return err
}

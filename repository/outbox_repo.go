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
	rows, err := tx.Query(`UPDATE outbox SET status = 'processing', processing_at = NOW()
WHERE id IN (
  SELECT id FROM outbox
  WHERE processed_at IS NULL AND status = 'pending'
  ORDER BY id
  LIMIT $1
  FOR UPDATE SKIP LOCKED
)
RETURNING id, type, caused_by, received_by, created_at, status`, limit)
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
	_, err := repo.db.Exec(`
		UPDATE outbox
		SET processed_at = NOW(), status = 'done', processing_at = NULL
		WHERE id = $1`, id)
	return err
}

// RecordEmailFailure increments attempts. Retries as pending until maxAttempts, then failed.
func (repo *OutboxRepository) RecordEmailFailure(id int64, maxAttempts int) error {
	_, err := repo.db.Exec(`
		UPDATE outbox
		SET attempts = attempts + 1,
		    processing_at = NULL,
		    status = CASE
		        WHEN attempts + 1 >= $2 THEN 'failed'
		        ELSE 'pending'
		    END
		WHERE id = $1`, id, maxAttempts)
	return err
}

// RequeueStuckProcessing moves long-running processing rows back to pending.
func (repo *OutboxRepository) RequeueStuckProcessing(olderThan time.Duration) (int64, error) {
	res, err := repo.db.Exec(`
		UPDATE outbox
		SET status = 'pending', processing_at = NULL
		WHERE status = 'processing'
		  AND processing_at IS NOT NULL
		  AND processing_at < NOW() - ($1 * INTERVAL '1 second')`,
		int64(olderThan.Seconds()),
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

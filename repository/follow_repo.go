package repository

import "database/sql"

type FollowRepository struct {
	db *sql.DB
}

func NewFollowRepository(db *sql.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) FollowUser(followingID, followerID int64, q DBTX) error {
	_, err := q.Exec(
		`INSERT INTO follows (follower_id, following_id) VALUES ($1, $2)`,
		followerID, followingID,
	)
	return err
}

func (r *FollowRepository) UnfollowUser(followingID, followerID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`,
		followerID, followingID,
	)
	return err
}

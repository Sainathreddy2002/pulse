package repository

import (
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

type LoginUser struct {
	ID             int64
	UserName       string
	HashedPassword string
}

type Me struct {
	ID        int64        `json:"id"`
	UserName  string       `json:"userName"`
	Email     string       `json:"email"`
	CreatedAt sql.NullTime `json:"createdAt"`
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(userName, email, passwordHash string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (user_name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		userName, email, passwordHash,
	).Scan(&id)
	return id, err
}

func (r *UserRepository) LoginUser(email string) (LoginUser, error) {
	var user LoginUser
	err := r.db.QueryRow(
		`SELECT id, user_name, password_hash FROM users WHERE email = $1 AND deleted_at IS NULL`,
		email,
	).Scan(&user.ID, &user.UserName, &user.HashedPassword)
	if err != nil {
		return LoginUser{}, err
	}
	return user, nil
}

func (r *UserRepository) Me(id int64) (Me, error) {
	var user Me
	err := r.db.QueryRow(
		`SELECT id, user_name, email, created_at FROM users WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&user.ID, &user.UserName, &user.Email, &user.CreatedAt)
	if err != nil {
		return Me{}, err
	}
	return user, nil
}

func (r *UserRepository) UserExists(id int64) (bool, error) {
	var exists int
	err := r.db.QueryRow(
		`SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

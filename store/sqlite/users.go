package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/rLukoyanov/w/core/models"
)

type UsersRepository struct {
	db *sql.DB
}

func (r *UsersRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UsersRepository) GetByID(id string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *UsersRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UsersRepository) GetByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE username = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *UsersRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UsersRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

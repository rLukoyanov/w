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
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password, role, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.Password, user.Role, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func scanUser(row interface{ Scan(...any) error }) (*models.User, error) {
	user := &models.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return user, nil
}

const userColumns = `id, username, email, password, role, created_at`

func (r *UsersRepository) GetByID(id string) (*models.User, error) {
	return scanUser(r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = ?`, id))
}

func (r *UsersRepository) GetByEmail(email string) (*models.User, error) {
	return scanUser(r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = ?`, email))
}

func (r *UsersRepository) GetByUsername(username string) (*models.User, error) {
	return scanUser(r.db.QueryRow(`SELECT `+userColumns+` FROM users WHERE username = ?`, username))
}

func (r *UsersRepository) GetAll() ([]*models.User, error) {
	rows, err := r.db.Query(`SELECT ` + userColumns + ` FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var result []*models.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}

func (r *UsersRepository) Update(user *models.User) error {
	result, err := r.db.Exec(
		`UPDATE users SET username = ?, email = ?, password = ?, role = ? WHERE id = ?`,
		user.Username, user.Email, user.Password, user.Role, user.ID,
	)
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
	result, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
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

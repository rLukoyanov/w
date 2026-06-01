package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/rLukoyanov/w/core/models"
)

type ServersRepository struct {
	db *sql.DB
}

func (r *ServersRepository) Create(server *models.Server) error {
	query := `
		INSERT INTO servers (id, name, owner_id, created_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, server.ID, server.Name, server.OwnerID, server.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

func (r *ServersRepository) GetByID(id string) (*models.Server, error) {
	query := `SELECT id, name, owner_id, created_at FROM servers WHERE id = ?`

	server := &models.Server{}
	err := r.db.QueryRow(query, id).Scan(
		&server.ID,
		&server.Name,
		&server.OwnerID,
		&server.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server by id: %w", err)
	}

	return server, nil
}

func (r *ServersRepository) GetByOwnerID(ownerID string) ([]*models.Server, error) {
	query := `SELECT id, name, owner_id, created_at FROM servers WHERE owner_id = ?`

	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get servers by owner id: %w", err)
	}
	defer rows.Close()

	var servers []*models.Server
	for rows.Next() {
		server := &models.Server{}
		err := rows.Scan(
			&server.ID,
			&server.Name,
			&server.OwnerID,
			&server.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan server: %w", err)
		}
		servers = append(servers, server)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return servers, nil
}

func (r *ServersRepository) Update(server *models.Server) error {
	query := `
		UPDATE servers 
		SET name = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, server.Name, server.ID)
	if err != nil {
		return fmt.Errorf("failed to update server: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("server not found")
	}

	return nil
}

func (r *ServersRepository) Delete(id string) error {
	query := `DELETE FROM servers WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete server: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("server not found")
	}

	return nil
}

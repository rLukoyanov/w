package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/rLukoyanov/w/core/models"
)

type ChannelsRepository struct {
	db *sql.DB
}

func (r *ChannelsRepository) Create(channel *models.Channel) error {
	query := `
		INSERT INTO channels (id, server_id, name, type, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, channel.ID, channel.ServerID, channel.Name, channel.Type, channel.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}

func (r *ChannelsRepository) GetByID(id string) (*models.Channel, error) {
	query := `SELECT id, server_id, name, type, created_at FROM channels WHERE id = ?`

	channel := &models.Channel{}
	err := r.db.QueryRow(query, id).Scan(
		&channel.ID,
		&channel.ServerID,
		&channel.Name,
		&channel.Type,
		&channel.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by id: %w", err)
	}

	return channel, nil
}

func (r *ChannelsRepository) GetByServerID(serverID string) ([]*models.Channel, error) {
	query := `SELECT id, server_id, name, type, created_at FROM channels WHERE server_id = ? ORDER BY created_at ASC`

	rows, err := r.db.Query(query, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels by server id: %w", err)
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		channel := &models.Channel{}
		err := rows.Scan(
			&channel.ID,
			&channel.ServerID,
			&channel.Name,
			&channel.Type,
			&channel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return channels, nil
}

func (r *ChannelsRepository) Update(channel *models.Channel) error {
	query := `
		UPDATE channels 
		SET name = ?, type = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, channel.Name, channel.Type, channel.ID)
	if err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

func (r *ChannelsRepository) Delete(id string) error {
	query := `DELETE FROM channels WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

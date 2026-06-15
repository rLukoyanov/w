package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rLukoyanov/w/core/models"
)

type MessagesRepository struct {
	db *sql.DB
}

func (r *MessagesRepository) Create(message *models.Message) error {
	query := `
		INSERT INTO messages (id, channel_id, author_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, message.ID, message.ChannelID, message.AuthorID, message.Content, message.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (r *MessagesRepository) GetByID(id string) (*models.Message, error) {
	query := `SELECT id, channel_id, author_id, content, created_at FROM messages WHERE id = ?`

	message := &models.Message{}
	err := r.db.QueryRow(query, id).Scan(
		&message.ID,
		&message.ChannelID,
		&message.AuthorID,
		&message.Content,
		&message.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message by id: %w", err)
	}

	return message, nil
}

func (r *MessagesRepository) GetByChannelID(channelID string, limit int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, channel_id, author_id, content, created_at 
		FROM messages 
		WHERE channel_id = ? 
		ORDER BY created_at DESC 
		LIMIT ?
	`

	rows, err := r.db.Query(query, channelID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by channel id: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.ChannelID,
			&message.AuthorID,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (r *MessagesRepository) GetByChannelIDBefore(channelID string, before time.Time, limit int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, channel_id, author_id, content, created_at 
		FROM messages 
		WHERE channel_id = ? AND created_at < ?
		ORDER BY created_at DESC 
		LIMIT ?
	`

	rows, err := r.db.Query(query, channelID, before, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages before: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(
			&message.ID,
			&message.ChannelID,
			&message.AuthorID,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (r *MessagesRepository) CountByChannelID(channelID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE channel_id = ?`, channelID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}
	return count, nil
}

func (r *MessagesRepository) Update(message *models.Message) error {
	query := `
		UPDATE messages 
		SET content = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, message.Content, message.ID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *MessagesRepository) Delete(id string) error {
	query := `DELETE FROM messages WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

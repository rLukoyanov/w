package sqlite

import (
	"database/sql"

	"github.com/rLukoyanov/w/core/models"
)

type AttachmentsRepository struct {
	db *sql.DB
}

func (r *AttachmentsRepository) Create(a *models.Attachment) error {
	_, err := r.db.Exec(
		`INSERT INTO attachments (id, channel_id, author_id, filename, size, content_type, storage_path, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		a.ID, a.ChannelID, a.AuthorID, a.Filename, a.Size, a.ContentType, a.StoragePath, a.CreatedAt,
	)
	return err
}

func (r *AttachmentsRepository) GetByID(id string) (*models.Attachment, error) {
	row := r.db.QueryRow(
		`SELECT id, channel_id, author_id, filename, size, content_type, storage_path, created_at
		 FROM attachments WHERE id = ?`, id,
	)
	var a models.Attachment
	err := row.Scan(&a.ID, &a.ChannelID, &a.AuthorID, &a.Filename, &a.Size, &a.ContentType, &a.StoragePath, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AttachmentsRepository) GetByChannelID(channelID string) ([]*models.Attachment, error) {
	rows, err := r.db.Query(
		`SELECT id, channel_id, author_id, filename, size, content_type, storage_path, created_at
		 FROM attachments WHERE channel_id = ? ORDER BY created_at ASC`, channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.ChannelID, &a.AuthorID, &a.Filename, &a.Size, &a.ContentType, &a.StoragePath, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *AttachmentsRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM attachments WHERE id = ?`, id)
	return err
}

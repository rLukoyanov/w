package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/rLukoyanov/w/core/models"
)

type ServerInvitesRepository struct {
	db *sql.DB
}

func (r *ServerInvitesRepository) Create(invite *models.ServerInvite) error {
	query := `INSERT INTO server_invites (id, server_id, created_by, code, max_uses, use_count, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, invite.ID, invite.ServerID, invite.CreatedBy, invite.Code, invite.MaxUses, invite.UseCount, invite.ExpiresAt, invite.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create invite: %w", err)
	}
	return nil
}

func (r *ServerInvitesRepository) GetByCode(code string) (*models.ServerInvite, error) {
	query := `SELECT id, server_id, created_by, code, max_uses, use_count, expires_at, created_at FROM server_invites WHERE code = ?`

	invite := &models.ServerInvite{}
	err := r.db.QueryRow(query, code).Scan(
		&invite.ID, &invite.ServerID, &invite.CreatedBy, &invite.Code,
		&invite.MaxUses, &invite.UseCount, &invite.ExpiresAt, &invite.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get invite by code: %w", err)
	}
	return invite, nil
}

func (r *ServerInvitesRepository) GetByServerID(serverID string) ([]*models.ServerInvite, error) {
	query := `SELECT id, server_id, created_by, code, max_uses, use_count, expires_at, created_at FROM server_invites WHERE server_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites by server: %w", err)
	}
	defer rows.Close()

	var invites []*models.ServerInvite
	for rows.Next() {
		invite := &models.ServerInvite{}
		if err := rows.Scan(&invite.ID, &invite.ServerID, &invite.CreatedBy, &invite.Code,
			&invite.MaxUses, &invite.UseCount, &invite.ExpiresAt, &invite.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan invite: %w", err)
		}
		invites = append(invites, invite)
	}

	if invites == nil {
		invites = []*models.ServerInvite{}
	}

	return invites, nil
}

func (r *ServerInvitesRepository) IncrementUseCount(id string) error {
	query := `UPDATE server_invites SET use_count = use_count + 1 WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to increment invite use count: %w", err)
	}
	return nil
}

func (r *ServerInvitesRepository) Delete(id string) error {
	query := `DELETE FROM server_invites WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}
	return nil
}

package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/rLukoyanov/w/core/models"
)

type ServerMembersRepository struct {
	db *sql.DB
}

func (r *ServerMembersRepository) Add(serverID, userID string) error {
	query := `INSERT INTO server_members (server_id, user_id, joined_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err := r.db.Exec(query, serverID, userID)
	if err != nil {
		return fmt.Errorf("failed to add server member: %w", err)
	}
	return nil
}

func (r *ServerMembersRepository) Remove(serverID, userID string) error {
	query := `DELETE FROM server_members WHERE server_id = ? AND user_id = ?`
	_, err := r.db.Exec(query, serverID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove server member: %w", err)
	}
	return nil
}

func (r *ServerMembersRepository) GetMembers(serverID string) ([]*models.ServerMember, error) {
	query := `SELECT server_id, user_id, joined_at FROM server_members WHERE server_id = ? ORDER BY joined_at`

	rows, err := r.db.Query(query, serverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server members: %w", err)
	}
	defer rows.Close()

	var members []*models.ServerMember
	for rows.Next() {
		m := &models.ServerMember{}
		if err := rows.Scan(&m.ServerID, &m.UserID, &m.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan server member: %w", err)
		}
		members = append(members, m)
	}

	if members == nil {
		members = []*models.ServerMember{}
	}

	return members, nil
}

func (r *ServerMembersRepository) IsMember(serverID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM server_members WHERE server_id = ? AND user_id = ?)`
	var exists bool
	err := r.db.QueryRow(query, serverID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check server membership: %w", err)
	}
	return exists, nil
}

func (r *ServerMembersRepository) GetServerIDsByUser(userID string) ([]string, error) {
	query := `SELECT server_id FROM server_members WHERE user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get server ids by user: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan server id: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

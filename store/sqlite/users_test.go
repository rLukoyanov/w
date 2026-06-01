package sqlite

import (
	"testing"
	"time"

	"github.com/rLukoyanov/w/core/models"
)

func TestUsersRepository_Create(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Apply migrations
	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	user := &models.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}

	err = store.Users().Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Verify user can be retrieved
	retrieved, err := store.Users().GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if retrieved == nil {
		t.Fatal("user not found")
	}

	if retrieved.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, retrieved.Username)
	}

	if retrieved.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, retrieved.Email)
	}
}

func TestUsersRepository_GetByEmail(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	user := &models.User{
		ID:        "user-456",
		Username:  "emailuser",
		Email:     "email@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}

	err = store.Users().Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	retrieved, err := store.Users().GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if retrieved == nil {
		t.Fatal("user not found")
	}

	if retrieved.ID != user.ID {
		t.Errorf("expected id %s, got %s", user.ID, retrieved.ID)
	}
}

func TestUsersRepository_GetByUsername(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	user := &models.User{
		ID:        "user-789",
		Username:  "uniqueuser",
		Email:     "unique@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}

	err = store.Users().Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	retrieved, err := store.Users().GetByUsername(user.Username)
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}

	if retrieved == nil {
		t.Fatal("user not found")
	}

	if retrieved.ID != user.ID {
		t.Errorf("expected id %s, got %s", user.ID, retrieved.ID)
	}
}

func TestUsersRepository_Update(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	user := &models.User{
		ID:        "user-update",
		Username:  "oldusername",
		Email:     "old@example.com",
		Password:  "oldpassword",
		CreatedAt: time.Now(),
	}

	err = store.Users().Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Update user
	user.Username = "newusername"
	user.Email = "new@example.com"
	user.Password = "newpassword"

	err = store.Users().Update(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	// Verify updates
	retrieved, err := store.Users().GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if retrieved.Username != "newusername" {
		t.Errorf("expected username newusername, got %s", retrieved.Username)
	}

	if retrieved.Email != "new@example.com" {
		t.Errorf("expected email new@example.com, got %s", retrieved.Email)
	}
}

func TestUsersRepository_Delete(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	user := &models.User{
		ID:        "user-delete",
		Username:  "deleteuser",
		Email:     "delete@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}

	err = store.Users().Create(user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Delete user
	err = store.Users().Delete(user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify user is gone
	retrieved, err := store.Users().GetByID(user.ID)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if retrieved != nil {
		t.Error("expected user to be deleted")
	}
}

func TestUsersRepository_GetByID_NotFound(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	retrieved, err := store.Users().GetByID("nonexistent")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrieved != nil {
		t.Error("expected nil for nonexistent user")
	}
}

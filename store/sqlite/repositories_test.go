package sqlite

import (
	"testing"
	"time"

	"github.com/rLukoyanov/w/core/models"
)

func TestServersRepository_Create(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	// Create a user first (for foreign key)
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

	server := &models.Server{
		ID:        "server-123",
		Name:      "Test Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}

	err = store.Servers().Create(server)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	retrieved, err := store.Servers().GetByID(server.ID)
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	if retrieved == nil {
		t.Fatal("server not found")
	}

	if retrieved.Name != server.Name {
		t.Errorf("expected name %s, got %s", server.Name, retrieved.Name)
	}
}

func TestServersRepository_GetByOwnerID(t *testing.T) {
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
		ID:        "user-owner",
		Username:  "owner",
		Email:     "owner@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	// Create multiple servers for same owner
	server1 := &models.Server{
		ID:        "server-1",
		Name:      "Server 1",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	server2 := &models.Server{
		ID:        "server-2",
		Name:      "Server 2",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}

	store.Servers().Create(server1)
	store.Servers().Create(server2)

	servers, err := store.Servers().GetByOwnerID(user.ID)
	if err != nil {
		t.Fatalf("failed to get servers by owner: %v", err)
	}

	if len(servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(servers))
	}
}

func TestServersRepository_Update(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-update",
		Name:      "Old Name",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	// Update server name
	server.Name = "New Name"
	err = store.Servers().Update(server)
	if err != nil {
		t.Fatalf("failed to update server: %v", err)
	}

	retrieved, err := store.Servers().GetByID(server.ID)
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	if retrieved.Name != "New Name" {
		t.Errorf("expected name 'New Name', got %s", retrieved.Name)
	}
}

func TestServersRepository_Delete(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-delete",
		Name:      "To Delete",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	err = store.Servers().Delete(server.ID)
	if err != nil {
		t.Fatalf("failed to delete server: %v", err)
	}

	retrieved, err := store.Servers().GetByID(server.ID)
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}

	if retrieved != nil {
		t.Error("expected server to be deleted")
	}
}

func TestChannelsRepository_Create(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	// Create user and server first
	user := &models.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-123",
		Name:      "Test Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-123",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}

	err = store.Channels().Create(channel)
	if err != nil {
		t.Fatalf("failed to create channel: %v", err)
	}

	retrieved, err := store.Channels().GetByID(channel.ID)
	if err != nil {
		t.Fatalf("failed to get channel: %v", err)
	}

	if retrieved == nil {
		t.Fatal("channel not found")
	}

	if retrieved.Name != channel.Name {
		t.Errorf("expected name %s, got %s", channel.Name, retrieved.Name)
	}
}

func TestChannelsRepository_GetByServerID(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	// Create multiple channels
	channel1 := &models.Channel{
		ID:        "channel-1",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	channel2 := &models.Channel{
		ID:        "channel-2",
		ServerID:  server.ID,
		Name:      "random",
		Type:      "text",
		CreatedAt: time.Now(),
	}

	store.Channels().Create(channel1)
	store.Channels().Create(channel2)

	channels, err := store.Channels().GetByServerID(server.ID)
	if err != nil {
		t.Fatalf("failed to get channels: %v", err)
	}

	if len(channels) != 2 {
		t.Errorf("expected 2 channels, got %d", len(channels))
	}
}

func TestChannelsRepository_Update(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-update",
		ServerID:  server.ID,
		Name:      "old-name",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	// Update channel
	channel.Name = "new-name"
	channel.Type = "voice"
	err = store.Channels().Update(channel)
	if err != nil {
		t.Fatalf("failed to update channel: %v", err)
	}

	retrieved, err := store.Channels().GetByID(channel.ID)
	if err != nil {
		t.Fatalf("failed to get channel: %v", err)
	}

	if retrieved.Name != "new-name" {
		t.Errorf("expected name 'new-name', got %s", retrieved.Name)
	}

	if retrieved.Type != "voice" {
		t.Errorf("expected type 'voice', got %s", retrieved.Type)
	}
}

func TestChannelsRepository_Delete(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-delete",
		ServerID:  server.ID,
		Name:      "to-delete",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	err = store.Channels().Delete(channel.ID)
	if err != nil {
		t.Fatalf("failed to delete channel: %v", err)
	}

	retrieved, err := store.Channels().GetByID(channel.ID)
	if err != nil {
		t.Fatalf("failed to get channel: %v", err)
	}

	if retrieved != nil {
		t.Error("expected channel to be deleted")
	}
}

func TestMessagesRepository_Create(t *testing.T) {
	store, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Migrate("./migrations")
	if err != nil {
		t.Fatalf("failed to apply migrations: %v", err)
	}

	// Create user, server, and channel first
	user := &models.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-123",
		Name:      "Test Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-123",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	message := &models.Message{
		ID:        "message-123",
		ChannelID: channel.ID,
		AuthorID:  user.ID,
		Content:   "Hello, World!",
		CreatedAt: time.Now(),
	}

	err = store.Messages().Create(message)
	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	retrieved, err := store.Messages().GetByID(message.ID)
	if err != nil {
		t.Fatalf("failed to get message: %v", err)
	}

	if retrieved == nil {
		t.Fatal("message not found")
	}

	if retrieved.Content != message.Content {
		t.Errorf("expected content %s, got %s", message.Content, retrieved.Content)
	}
}

func TestMessagesRepository_GetByChannelID(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-1",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	// Create multiple messages
	for i := 1; i <= 5; i++ {
		msg := &models.Message{
			ID:        string(rune('a' + i)),
			ChannelID: channel.ID,
			AuthorID:  user.ID,
			Content:   "Message " + string(rune('0'+i)),
			CreatedAt: time.Now(),
		}
		store.Messages().Create(msg)
		time.Sleep(time.Millisecond) // Ensure different timestamps
	}

	messages, err := store.Messages().GetByChannelID(channel.ID, 3)
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(messages) != 3 {
		t.Errorf("expected 3 messages with limit, got %d", len(messages))
	}

	// Test without limit (should use default 50)
	allMessages, err := store.Messages().GetByChannelID(channel.ID, 0)
	if err != nil {
		t.Fatalf("failed to get all messages: %v", err)
	}

	if len(allMessages) != 5 {
		t.Errorf("expected 5 messages without limit, got %d", len(allMessages))
	}
}

func TestMessagesRepository_Update(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-1",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	message := &models.Message{
		ID:        "message-update",
		ChannelID: channel.ID,
		AuthorID:  user.ID,
		Content:   "Original content",
		CreatedAt: time.Now(),
	}
	store.Messages().Create(message)

	// Update message
	message.Content = "Edited content"
	err = store.Messages().Update(message)
	if err != nil {
		t.Fatalf("failed to update message: %v", err)
	}

	retrieved, err := store.Messages().GetByID(message.ID)
	if err != nil {
		t.Fatalf("failed to get message: %v", err)
	}

	if retrieved.Content != "Edited content" {
		t.Errorf("expected content 'Edited content', got %s", retrieved.Content)
	}
}

func TestMessagesRepository_Delete(t *testing.T) {
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
		ID:        "user-1",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
	}
	store.Users().Create(user)

	server := &models.Server{
		ID:        "server-1",
		Name:      "Server",
		OwnerID:   user.ID,
		CreatedAt: time.Now(),
	}
	store.Servers().Create(server)

	channel := &models.Channel{
		ID:        "channel-1",
		ServerID:  server.ID,
		Name:      "general",
		Type:      "text",
		CreatedAt: time.Now(),
	}
	store.Channels().Create(channel)

	message := &models.Message{
		ID:        "message-delete",
		ChannelID: channel.ID,
		AuthorID:  user.ID,
		Content:   "To be deleted",
		CreatedAt: time.Now(),
	}
	store.Messages().Create(message)

	err = store.Messages().Delete(message.ID)
	if err != nil {
		t.Fatalf("failed to delete message: %v", err)
	}

	retrieved, err := store.Messages().GetByID(message.ID)
	if err != nil {
		t.Fatalf("failed to get message: %v", err)
	}

	if retrieved != nil {
		t.Error("expected message to be deleted")
	}
}

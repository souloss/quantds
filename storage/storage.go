// Package storage provides optional data persistence for quantds.
// It defines the Storage interface and provides file-based and in-memory implementations.
// Storage is not enabled by default; users can inject it via Manager options.
package storage

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Storage defines the interface for optional data persistence.
// Implementations must be safe for concurrent use.
type Storage interface {
	// Save stores data with the given key and TTL.
	// If TTL is zero, the data does not expire.
	Save(ctx context.Context, key string, data []byte, ttl time.Duration) error
	// Load retrieves data by key. Returns ErrNotFound if the key does not exist or has expired.
	Load(ctx context.Context, key string) ([]byte, error)
	// Delete removes data by key.
	Delete(ctx context.Context, key string) error
	// Close releases any resources held by the storage.
	Close() error
}

// ErrNotFound is returned when a key is not found in storage.
type ErrNotFound struct {
	Key string
}

func (e *ErrNotFound) Error() string {
	return "storage: key not found: " + e.Key
}

// MemoryStorage is an in-memory implementation of Storage.
// It is useful for testing and as a no-persistence default.
type MemoryStorage struct {
	mu    sync.RWMutex
	items map[string]*memoryItem
}

type memoryItem struct {
	data      []byte
	expiresAt time.Time
}

// NewMemoryStorage creates a new in-memory storage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		items: make(map[string]*memoryItem),
	}
}

func (m *MemoryStorage) Save(_ context.Context, key string, data []byte, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := &memoryItem{data: data}
	if ttl > 0 {
		item.expiresAt = time.Now().Add(ttl)
	}
	m.items[key] = item
	return nil
}

func (m *MemoryStorage) Load(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return nil, &ErrNotFound{Key: key}
	}
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return nil, &ErrNotFound{Key: key}
	}
	// Return a copy to prevent mutation
	result := make([]byte, len(item.data))
	copy(result, item.data)
	return result, nil
}

func (m *MemoryStorage) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
	return nil
}

func (m *MemoryStorage) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = nil
	return nil
}

// FileStorage is a file-system-based implementation of Storage.
// Each key is stored as a separate file in the specified directory.
type FileStorage struct {
	dir string
}

type fileEntry struct {
	Data      []byte `json:"data"`
	ExpiresAt int64  `json:"expires_at,omitempty"` // UnixMilli timestamp, 0 = no expiry
}

// NewFileStorage creates a new file-based storage in the given directory.
// The directory is created if it does not exist.
func NewFileStorage(dir string) (*FileStorage, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileStorage{dir: dir}, nil
}

func (f *FileStorage) Save(_ context.Context, key string, data []byte, ttl time.Duration) error {
	entry := fileEntry{Data: data}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl).UnixMilli()
	}

	path := f.filePath(key)
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func (f *FileStorage) Load(_ context.Context, key string) ([]byte, error) {
	path := f.filePath(key)
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ErrNotFound{Key: key}
		}
		return nil, err
	}

	var entry fileEntry
	if err := json.Unmarshal(b, &entry); err != nil {
		return nil, err
	}

	if entry.ExpiresAt > 0 && time.Now().UnixMilli() > entry.ExpiresAt {
		os.Remove(path)
		return nil, &ErrNotFound{Key: key}
	}

	return entry.Data, nil
}

func (f *FileStorage) Delete(_ context.Context, key string) error {
	path := f.filePath(key)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (f *FileStorage) Close() error {
	return nil
}

func (f *FileStorage) filePath(key string) string {
	return filepath.Join(f.dir, key+".json")
}

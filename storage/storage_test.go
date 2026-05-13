package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMemoryStorage_SaveLoad(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	ctx := context.Background()

	// Save and load
	err := s.Save(ctx, "key1", []byte("value1"), 0)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := s.Load(ctx, "key1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if string(data) != "value1" {
		t.Errorf("Load = %q, want %q", data, "value1")
	}
}

func TestMemoryStorage_NotFound(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	_, err := s.Load(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
	if _, ok := err.(*ErrNotFound); !ok {
		t.Errorf("expected *ErrNotFound, got %T", err)
	}
}

func TestMemoryStorage_TTL(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	ctx := context.Background()

	// Save with short TTL
	err := s.Save(ctx, "key1", []byte("value1"), 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Should be available immediately
	_, err = s.Load(ctx, "key1")
	if err != nil {
		t.Errorf("Load before expiry: %v", err)
	}

	// Should be expired after TTL
	time.Sleep(100 * time.Millisecond)
	_, err = s.Load(ctx, "key1")
	if err == nil {
		t.Error("expected error after TTL expiry")
	}
}

func TestMemoryStorage_Delete(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	ctx := context.Background()

	s.Save(ctx, "key1", []byte("value1"), 0)
	s.Delete(ctx, "key1")

	_, err := s.Load(ctx, "key1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestMemoryStorage_Overwrite(t *testing.T) {
	s := NewMemoryStorage()
	defer s.Close()

	ctx := context.Background()

	s.Save(ctx, "key1", []byte("value1"), 0)
	s.Save(ctx, "key1", []byte("value2"), 0)

	data, err := s.Load(ctx, "key1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if string(data) != "value2" {
		t.Errorf("Load = %q, want %q", data, "value2")
	}
}

func TestFileStorage_SaveLoad(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage: %v", err)
	}
	defer s.Close()

	ctx := context.Background()

	err = s.Save(ctx, "key1", []byte("value1"), 0)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	data, err := s.Load(ctx, "key1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if string(data) != "value1" {
		t.Errorf("Load = %q, want %q", data, "value1")
	}
}

func TestFileStorage_NotFound(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage: %v", err)
	}
	defer s.Close()

	_, err = s.Load(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent key")
	}
	if _, ok := err.(*ErrNotFound); !ok {
		t.Errorf("expected *ErrNotFound, got %T", err)
	}
}

func TestFileStorage_TTL(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage: %v", err)
	}
	defer s.Close()

	ctx := context.Background()

	err = s.Save(ctx, "key1", []byte("value1"), 50*time.Millisecond)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Should be available immediately
	_, err = s.Load(ctx, "key1")
	if err != nil {
		t.Errorf("Load before expiry: %v", err)
	}

	// Should be expired after TTL
	time.Sleep(100 * time.Millisecond)
	_, err = s.Load(ctx, "key1")
	if err == nil {
		t.Error("expected error after TTL expiry")
	}
}

func TestFileStorage_Delete(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage: %v", err)
	}
	defer s.Close()

	ctx := context.Background()

	s.Save(ctx, "key1", []byte("value1"), 0)
	s.Delete(ctx, "key1")

	_, err = s.Load(ctx, "key1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestFileStorage_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "subdir", "data")
	s, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("NewFileStorage: %v", err)
	}
	defer s.Close()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("NewFileStorage should create directory")
	}
}

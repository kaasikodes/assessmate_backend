package payment

import (
	"errors"
	"fmt"
	"strings"
)

type Meta map[string]string

func NewMeta() (*Meta, error) {
	meta := make(Meta)
	return &meta, nil

}

// Add inserts a key-value pair if the key does not already exist.
func (m *Meta) Add(key, value string) error {
	if m == nil {
		return errors.New("meta is nil")
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("key cannot be empty")
	}

	if _, exists := (*m)[key]; exists {
		return fmt.Errorf("key %q already exists in meta", key)
	}

	(*m)[key] = value
	return nil
}

// Get returns the value associated with a key.
func (m *Meta) Get(key string) (string, bool) {
	if m == nil {
		return "", false
	}

	value, exists := (*m)[key]
	return value, exists
}

// Remove deletes a key from the map if it exists.
func (m *Meta) Remove(key string) error {
	if m == nil {
		return errors.New("meta is nil")
	}

	key = strings.TrimSpace(key)
	if _, exists := (*m)[key]; !exists {
		return fmt.Errorf("key %q does not exist in meta", key)
	}

	delete(*m, key)
	return nil
}

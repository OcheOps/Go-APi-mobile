package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Note represents a field note
type Note struct {
	// Unique ID of the note
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Title of the note
	Title string `json:"title" example:"My Field Note"`
	// Body content of the note
	Body string `json:"body" example:"Observed some interesting patterns..."`
	// List of tags associated with the note
	Tags Tags `json:"tags" swaggertype:"array,string" example:"nature,assessment"`
	// Last update timestamp in RFC3339 format
	UpdatedAt time.Time `json:"updatedAt" example:"2026-01-27T20:20:45+01:00"`
	// Creation timestamp in RFC3339 format
	CreatedAt time.Time `json:"createdAt" example:"2026-01-27T20:20:45+01:00"`
}

// CreateNoteRequest body for creating a note
type CreateNoteRequest struct {
	// Title of the note (required)
	Title string `json:"title" example:"New Discovery"`
	// Body of the note
	Body string `json:"body" example:"Detailed description..."`
	// Tags for the note
	Tags []string `json:"tags" example:"science,important"`
}

// UpdateNoteRequest body for updating a note
type UpdateNoteRequest struct {
	// New title (optional)
	Title *string `json:"title" example:"Updated Discovery"`
	// New body (optional)
	Body *string `json:"body" example:"Updated description..."`
	// New tags (optional)
	Tags *[]string `json:"tags" example:"science,update"`
	// Timestamp of the note version being updated (required for conflict detection)
	UpdatedAt time.Time `json:"updatedAt" example:"2026-01-27T20:20:45+01:00"`
}

// Tags handling for SQLite text storage
type Tags []string

func (t Tags) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "[]", nil
	}
	return json.Marshal(t)
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = []string{}
		return nil
	}
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), t)
	case []byte:
		return json.Unmarshal(v, t)
	}
	return errors.New("failed to scan tags")
}

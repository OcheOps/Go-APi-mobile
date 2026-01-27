package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Helper for JSON responses
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

// ListNotes godoc
// @Summary      List all notes
// @Description  Get all notes, ordered by last update time descending
// @Tags         notes
// @Produce      json
// @Success      200  {array}   Note
// @Failure      500  {object}  map[string]string
// @Router       /notes [get]
func ListNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, body, tags, updated_at, created_at FROM notes ORDER BY updated_at DESC")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		var up, cr string
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.Tags, &up, &cr); err != nil {
			errorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		n.UpdatedAt, _ = time.Parse(time.RFC3339, up)
		n.CreatedAt, _ = time.Parse(time.RFC3339, cr)
		notes = append(notes, n)
	}
	jsonResponse(w, http.StatusOK, notes)
}

// CreateNote godoc
// @Summary      Create a new note
// @Description  Create a new note with the provided title, body, and tags
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        note  body      CreateNoteRequest  true  "Note content"
// @Success      201   {object}  Note
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notes [post]
func CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		errorResponse(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Tags == nil {
		req.Tags = []string{}
	}

	now := time.Now().UTC()
	note := Note{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Body:      req.Body,
		Tags:      req.Tags,
		UpdatedAt: now,
		CreatedAt: now,
	}

	query := "INSERT INTO notes (id, title, body, tags, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query,
		note.ID,
		note.Title,
		note.Body,
		note.Tags,
		note.UpdatedAt.Format(time.RFC3339),
		note.CreatedAt.Format(time.RFC3339),
	)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, note)
}

// UpdateNote godoc
// @Summary      Update a note
// @Description  Update specific fields of a note. Requires valid 'updatedAt' for conflict detection.
// @Tags         notes
// @Accept       json
// @Produce      json
// @Param        id    path      string             true  "Note ID"
// @Param        note  body      UpdateNoteRequest  true  "Update content"
// @Success      200   {object}  Note
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      409   {object}  Note "Conflict: Returns server version of the note"
// @Failure      500   {object}  map[string]string
// @Router       /notes/{id} [patch]
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UpdatedAt.IsZero() {
		errorResponse(w, http.StatusBadRequest, "updatedAt is required")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer tx.Rollback()

	var current Note
	var upStr, crStr string

	err = tx.QueryRow("SELECT id, title, body, tags, updated_at, created_at FROM notes WHERE id = ?", id).Scan(
		&current.ID, &current.Title, &current.Body, &current.Tags, &upStr, &crStr,
	)

	if err == sql.ErrNoRows {
		errorResponse(w, http.StatusNotFound, "Note not found")
		return
	} else if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	current.UpdatedAt, _ = time.Parse(time.RFC3339, upStr)
	current.CreatedAt, _ = time.Parse(time.RFC3339, crStr)

	// Conflict Detection
	if req.UpdatedAt.Before(current.UpdatedAt) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409
		json.NewEncoder(w).Encode(current)
		return
	}

	// Prepare updates
	newTitle := current.Title
	if req.Title != nil {
		newTitle = *req.Title
	}
	newBody := current.Body
	if req.Body != nil {
		newBody = *req.Body
	}
	newTags := current.Tags
	if req.Tags != nil {
		newTags = *req.Tags
	}
	newUpdatedAt := time.Now().UTC()

	updateQuery := "UPDATE notes SET title = ?, body = ?, tags = ?, updated_at = ? WHERE id = ?"
	_, err = tx.Exec(updateQuery, newTitle, newBody, newTags, newUpdatedAt.Format(time.RFC3339), id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	current.Title = newTitle
	current.Body = newBody
	current.Tags = newTags
	current.UpdatedAt = newUpdatedAt

	jsonResponse(w, http.StatusOK, current)
}

// DeleteNote godoc
// @Summary      Delete a note
// @Description  Delete a note by its ID
// @Tags         notes
// @Produce      json
// @Param        id   path      string  true  "Note ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notes/{id} [delete]
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		errorResponse(w, http.StatusNotFound, "Note not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

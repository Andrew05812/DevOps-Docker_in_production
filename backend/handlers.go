package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := "ok"

	if a.pool != nil {
		if err := a.pool.Ping(r.Context()); err != nil {
			status = "degraded"
		}
	}

	code := http.StatusOK
	if status != "ok" {
		code = http.StatusServiceUnavailable
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": status,
	})
}

func (a *App) getNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := a.pool.Query(r.Context(), `
		SELECT id, title, content, created_at
		FROM notes
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := make([]Note, 0)

	for rows.Next() {
		var n Note

		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
		); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		notes = append(notes, n)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notes)
}

func (a *App) createNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)

	if req.Title == "" || req.Content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	var note Note

	err := a.pool.QueryRow(
		r.Context(),
		`
		INSERT INTO notes (title, content)
		VALUES ($1, $2)
		RETURNING id, title, content, created_at
		`,
		req.Title,
		req.Content,
	).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)

	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(note)
}

func (a *App) getNote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/notes/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var note Note

	err = a.pool.QueryRow(
		r.Context(),
		`
		SELECT id, title, content, created_at
		FROM notes
		WHERE id = $1
		`,
		id,
	).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
	)

	if err != nil {
		http.Error(w, "note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(note)
}

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vizvim/board/validator"
)

var ErrBoardNotFound = errors.New("board not found")

type Board struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Image     []byte    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Version   int       `json:"version"`
}

func (b Board) Validate() map[string]string {
	v := validator.New()

	v.Check(b.Name != "", "name", "must be provided")
	v.Check(b.Image != nil, "image", "must be provided")

	if v.Valid() {
		return nil
	}

	return v.Errors
}

func (d *DB) CreateBoard(ctx context.Context, b *Board) error {
	query := `
    INSERT INTO boards (name, image) 
    VALUES ($1, $2)
    RETURNING id, created_at, updated_at, version`

	args := []any{b.Name, b.Image}

	err := d.QueryRowContext(ctx, query, args...).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt, &b.Version)
	if err != nil {
		return fmt.Errorf("error creating board: %v", err)
	}

	return nil
}

func (d *DB) GetBoard(id uuid.UUID) (Board, error) {
	query := `SELECT id, name, image, created_at, updated_at, version FROM boards WHERE id = $1`

	var board Board

	err := d.QueryRow(query, id).Scan(&board.ID, &board.Name, &board.Image, &board.CreatedAt, &board.UpdatedAt, &board.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return Board{}, ErrBoardNotFound
		}

		return Board{}, fmt.Errorf("error getting board: %v", err)
	}

	return board, nil
}

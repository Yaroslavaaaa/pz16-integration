package repo

import (
	"context"
	"database/sql"
	"errors"

	"example.com/pz16/internal/models"
)

type NoteRepo struct{ DB *sql.DB }

func (r NoteRepo) Create(ctx context.Context, n *models.Note) error {
	return r.DB.QueryRowContext(ctx,
		`INSERT INTO notes(title, content) VALUES($1,$2) RETURNING id`,
		n.Title, n.Content,
	).Scan(&n.ID)
}

func (r NoteRepo) Get(ctx context.Context, id int64) (models.Note, error) {
	var n models.Note
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, title, content, created_at, updated_at FROM notes WHERE id=$1`, id,
	).Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.Note{}, errors.New("not found")
	}
	return n, err
}

func (r NoteRepo) Update(ctx context.Context, n *models.Note) error {
	res, err := r.DB.ExecContext(ctx,
		`UPDATE notes SET title=$1, content=$2 WHERE id=$3`,
		n.Title, n.Content, n.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r NoteRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM notes WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r NoteRepo) List(ctx context.Context) ([]models.Note, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, title, content, created_at, updated_at FROM notes ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

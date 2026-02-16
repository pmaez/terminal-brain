package core

import (
	"context"
	"database/sql"
	"fmt"
)

func NewSQLiteStores(db *sql.DB) Stores {
	return Stores{
		Stack:   &stackStore{db: db},
		Dump:    &dumpStore{db: db},
		Duck:    &duckStore{db: db},
		Command: &commandStore{db: db},
		Time:    &timeStore{db: db},
	}
}

type stackStore struct {
	db *sql.DB
}

type dumpStore struct {
	db *sql.DB
}

type duckStore struct {
	db *sql.DB
}

type commandStore struct {
	db *sql.DB
}

type timeStore struct {
	db *sql.DB
}

func (s *stackStore) List(ctx context.Context) ([]StackTask, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, title, position, created_at, updated_at FROM stack_tasks ORDER BY position ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []StackTask
	for rows.Next() {
		var item StackTask
		if err := rows.Scan(&item.ID, &item.Title, &item.Position, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *stackStore) Add(ctx context.Context, title string) (StackTask, error) {
	var nextPos int
	if err := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(position), -1) + 1 FROM stack_tasks`).Scan(&nextPos); err != nil {
		return StackTask{}, err
	}

	now := nowString()
	res, err := s.db.ExecContext(ctx, `INSERT INTO stack_tasks (title, position, created_at, updated_at) VALUES (?, ?, ?, ?)`, title, nextPos, now, now)
	if err != nil {
		return StackTask{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return StackTask{}, err
	}
	return StackTask{ID: id, Title: title, Position: nextPos, CreatedAt: now, UpdatedAt: now}, nil
}

func (s *stackStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM stack_tasks WHERE id = ?`, id)
	return err
}

func (s *stackStore) UpdateOrder(ctx context.Context, ids []int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for idx, id := range ids {
		if _, err := tx.ExecContext(ctx, `UPDATE stack_tasks SET position = ?, updated_at = ? WHERE id = ?`, idx, nowString(), id); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *dumpStore) List(ctx context.Context, limit int) ([]DumpEntry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, body, created_at FROM dump_entries ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DumpEntry
	for rows.Next() {
		var item DumpEntry
		if err := rows.Scan(&item.ID, &item.Body, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *dumpStore) Add(ctx context.Context, body string) (DumpEntry, error) {
	now := nowString()
	res, err := s.db.ExecContext(ctx, `INSERT INTO dump_entries (body, created_at) VALUES (?, ?)`, body, now)
	if err != nil {
		return DumpEntry{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return DumpEntry{}, err
	}
	return DumpEntry{ID: id, Body: body, CreatedAt: now}, nil
}

func (s *duckStore) List(ctx context.Context, limit int) ([]DuckSession, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, q1, q2, q3, q4, created_at FROM duck_sessions ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []DuckSession
	for rows.Next() {
		var item DuckSession
		if err := rows.Scan(&item.ID, &item.Q1, &item.Q2, &item.Q3, &item.Q4, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *duckStore) Add(ctx context.Context, q1, q2, q3, q4 string) (DuckSession, error) {
	now := nowString()
	res, err := s.db.ExecContext(ctx, `INSERT INTO duck_sessions (q1, q2, q3, q4, created_at) VALUES (?, ?, ?, ?, ?)`, q1, q2, q3, q4, now)
	if err != nil {
		return DuckSession{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return DuckSession{}, err
	}
	return DuckSession{ID: id, Q1: q1, Q2: q2, Q3: q3, Q4: q4, CreatedAt: now}, nil
}

func (s *duckStore) Update(ctx context.Context, id int64, q1, q2, q3, q4 string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE duck_sessions SET q1 = ?, q2 = ?, q3 = ?, q4 = ? WHERE id = ?`, q1, q2, q3, q4, id)
	return err
}

func (s *duckStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM duck_sessions WHERE id = ?`, id)
	return err
}

func (s *commandStore) Search(ctx context.Context, query string, limit int) ([]CommandEntry, error) {
	like := fmt.Sprintf("%%%s%%", query)
	rows, err := s.db.QueryContext(ctx, `SELECT id, command, description, uses, last_used, created_at FROM commands WHERE command LIKE ? OR description LIKE ? ORDER BY last_used DESC, uses DESC, id DESC LIMIT ?`, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CommandEntry
	for rows.Next() {
		var item CommandEntry
		if err := rows.Scan(&item.ID, &item.Command, &item.Description, &item.Uses, &item.LastUsed, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *commandStore) Add(ctx context.Context, command, description string) (CommandEntry, error) {
	now := nowString()
	res, err := s.db.ExecContext(ctx, `INSERT INTO commands (command, description, uses, last_used, created_at) VALUES (?, ?, 0, '', ?)`, command, description, now)
	if err != nil {
		return CommandEntry{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return CommandEntry{}, err
	}
	return CommandEntry{ID: id, Command: command, Description: description, Uses: 0, LastUsed: "", CreatedAt: now}, nil
}

func (s *commandStore) Touch(ctx context.Context, id int64) error {
	now := nowString()
	_, err := s.db.ExecContext(ctx, `UPDATE commands SET uses = uses + 1, last_used = ? WHERE id = ?`, now, id)
	return err
}

func (s *timeStore) List(ctx context.Context, limit int) ([]TimeLog, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, label, start, end, created_at FROM time_logs ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []TimeLog
	for rows.Next() {
		var item TimeLog
		if err := rows.Scan(&item.ID, &item.Label, &item.Start, &item.End, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *timeStore) Start(ctx context.Context, label, start string) (TimeLog, error) {
	created := nowString()
	res, err := s.db.ExecContext(ctx, `INSERT INTO time_logs (label, start, end, created_at) VALUES (?, ?, '', ?)`, label, start, created)
	if err != nil {
		return TimeLog{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return TimeLog{}, err
	}
	return TimeLog{ID: id, Label: label, Start: start, End: "", CreatedAt: created}, nil
}

func (s *timeStore) Stop(ctx context.Context, id int64, end string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE time_logs SET end = ? WHERE id = ?`, end, id)
	return err
}

func (s *timeStore) Delete(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM time_logs WHERE id = ?`, id)
	return err
}

func (s *timeStore) Clear(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM time_logs`)
	return err
}

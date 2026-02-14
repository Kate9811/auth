package auth

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/Denis/project_auth/internal/client/db"
	"github.com/Denis/project_auth/internal/model"
	"github.com/Denis/project_auth/internal/repository"
)

const (
	tableName    = "auth"
	colID        = "id"
	colName      = "name"
	colEmail     = "email"
	colPassword  = "password_hash"
	colRole      = "role"
	colCreatedAt = "created_at"
	colUpdatedAt = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(colName, colEmail, colPassword, colRole).
		Values(info.Name, info.Email, info.PasswordHash, info.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.Create", // Исправлено имя
		QueryRaw: query,
	}

	var id int64
	// Правильно: через r.db.DB().QueryRowContext()
	err = r.db.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, fmt.Errorf("exec query: %w", err)
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Auth, error) {
	builder := sq.Select(
		colID, colName, colEmail, colPassword, colRole, colCreatedAt, colUpdatedAt,
	).
		From(tableName).
		Where(sq.Eq{colID: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.Get", // Исправлено имя
		QueryRaw: query,
	}

	var auth model.Auth
	// Правильно: через r.db.DB().QueryRowContext()
	err = r.db.DB().ScanOneContext(ctx, &auth, q, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return &auth, nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.AuthInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{colID: id})

	// Добавляем только те поля, которые не пустые
	if info.Name != "" {
		builder = builder.Set(colName, info.Name)
	}

	if info.Email != "" {
		builder = builder.Set(colEmail, info.Email)
	}

	if info.PasswordHash != "" {
		builder = builder.Set(colPassword, info.PasswordHash)
	}

	if info.Role != "" {
		builder = builder.Set(colRole, info.Role)
	}

	// Всегда обновляем updated_at
	builder = builder.Set(colUpdatedAt, sq.Expr("NOW()"))

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	q := db.Query{ // Добавьте db.Query!
		Name:     "auth_repository.Update",
		QueryRaw: query,
	}

	// ❌ НЕПРАВИЛЬНО: r.db.ExecContext(ctx, query, args...)
	// ✅ ПРАВИЛЬНО: через r.db.DB().ExecContext()
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{colID: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	q := db.Query{ // Добавьте db.Query!
		Name:     "auth_repository.Delete",
		QueryRaw: query,
	}

	// ❌ НЕПРАВИЛЬНО: r.db.ExecContext(ctx, query, args...)
	// ✅ ПРАВИЛЬНО: через r.db.DB().ExecContext()
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

package auth

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"

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
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
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
	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
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

	var auth model.Auth
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&auth.ID,
		&auth.Info.Name,
		&auth.Info.Email,
		&auth.Info.PasswordHash,
		&auth.Info.Role,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)
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

	// Добавьте лог для отладки!
	log.Printf("[SQL UPDATE] Query: %s, Args: %v", query, args)

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("[SQL UPDATE ERROR] %v", err)
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

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID        int64        `db:"id"`         // Первичный ключ, обычно автоинкремент
	Info      AuthInfo     `db:""`           // Игнорируется, т.к. это вложенная структура
	CreatedAt time.Time    `db:"created_at"` // Время создания записи
	UpdatedAt sql.NullTime `db:"updated_at"` // Время обновления (может быть NULL)
}

type AuthInfo struct {
	Name         string `db:"name"`          // Имя пользователя
	Email        string `db:"email"`         // Email (уникальный)
	PasswordHash string `db:"password_hash"` // Хеш пароля
	Role         string `db:"role"`          // Роль пользователя
}

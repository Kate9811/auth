package pg

import (
	"context"
	"fmt"
	"log"

	// scany - библиотека для маппинга результатов SQL в структуры Go
	"github.com/georgysavva/scany/pgxscan"
	// pgconn - низкоуровневые операции PostgreSQL (команды, результаты)
	"github.com/jackc/pgconn"
	// pgx - драйвер PostgreSQL для Go
	"github.com/jackc/pgx/v4"
	// pgxpool - пул соединений для pgx
	"github.com/jackc/pgx/v4/pgxpool"

	// Внутренние пакеты проекта
	"github.com/Denis/project_auth/internal/client/db"
	"github.com/Denis/project_auth/internal/client/db/prettier"
)

// key - тип для ключей контекста (для type-safe доступа)
type key string

// TxKey - ключ для хранения транзакции в контексте
const (
	TxKey key = "tx"
)

// pg - структура, реализующая интерфейс db.DB
// Инкапсулирует пул соединений с PostgreSQL
type pg struct {
	dbc *pgxpool.Pool // Пул соединений с БД
}

// NewDB - конструктор, создает новую реализацию интерфейса db.DB
// Принимает готовый пул соединений pgxpool.Pool
func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc, // Сохраняем переданный пул соединений
	}
}

// ScanOneContext - выполняет SQL-запрос и маппит ОДНУ строку результата в dest
// Используется для SELECT запросов, возвращающих одну строку
func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...) // Логируем запрос для отладки

	// Выполняем запрос и получаем результат
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err // Возвращаем ошибку если запрос не удался
	}

	// Используем scany для маппинга строки в структуру dest
	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext - выполняет SQL-запрос и маппит ВСЕ строки результата в dest (слайс структур)
// Используется для SELECT запросов, возвращающих несколько строк
func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...) // Логируем запрос

	// Выполняем запрос и получаем все строки результата
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	// Используем scany для маппинга всех строк в слайс структур dest
	return pgxscan.ScanAll(dest, rows)
}

// ExecContext - выполняет SQL-запрос НЕ возвращающий строки (INSERT, UPDATE, DELETE)
// Возвращает CommandTag с информацией о количестве затронутых строк
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...) // Логируем запрос

	// Проверяем, есть ли транзакция в контексте
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		// Если есть транзакция - выполняем запрос в её рамках
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	// Если транзакции нет - выполняем запрос напрямую через пул
	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext - выполняет SQL-запрос и возвращает МНОЖЕСТВО строк результата
// Используется для SELECT запросов, которые нужно обрабатывать построчно
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {
	logQuery(ctx, q, args...) // Логируем запрос

	// Проверяем, есть ли транзакция в контексте
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		// Если есть транзакция - выполняем запрос в её рамках
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	// Если транзакции нет - выполняем запрос напрямую через пул
	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext - выполняет SQL-запрос и возвращает ОДНУ строку результата
// Оптимизирован для запросов, которые гарантированно возвращают одну строку
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...) // Логируем запрос

	// Проверяем, есть ли транзакция в контексте
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		// Если есть транзакция - выполняем запрос в её рамках
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	// Если транзакции нет - выполняем запрос напрямую через пул
	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// BeginTx - начинает новую транзакцию с указанными параметрами
// txOptions может содержать уровень изоляции, режим доступа (только чтение/запись)
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	// Делегируем создание транзакции пулу соединений
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping - проверяет, что соединение с БД активно
// Используется для health checks и readiness probes
func (p *pg) Ping(ctx context.Context) error {
	// Делегируем проверку пулу соединений
	return p.dbc.Ping(ctx)
}

// Close - закрывает все соединения в пуле
// Должен вызываться при завершении работы приложения
func (p *pg) Close() {
	p.dbc.Close()
}

// // MakeContextTx - создает новый контекст с привязанной транзакцией
// // Позволяет передавать транзакцию через контекст без явной передачи параметром
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	// Сохраняем транзакцию в контексте под ключом TxKey
	return context.WithValue(ctx, TxKey, tx)
}

// logQuery - вспомогательная функция для логирования SQL-запросов
// Форматирует запрос, подставляя значения аргументов вместо плейсхолдеров
func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	// Преобразуем запрос в читаемый вид, заменяя плейсхолдеры $1, $2 на фактические значения
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)

	// Логируем: контекст, имя запроса и сам запрос с подставленными значениями
	log.Println(
		ctx,                                   // Контекст (может содержать ID запроса, таймауты и т.д.)
		fmt.Sprintf("sql: %s", q.Name),        // Имя запроса для идентификации
		fmt.Sprintf("query: %s", prettyQuery), // Сам SQL-запрос с данными
	)
}

package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

// Константы для типов плейсхолдеров в SQL-запросах
const (
	PlaceholderDollar   = "$" // Используется в PostgreSQL: $1, $2, $3
	PlaceholderQuestion = "?" // Используется в MySQL/SQLite: ?, ?, ?
)

// Pretty - функция для форматирования SQL-запросов
// Заменяет плейсхолдеры ($1, $2 или ?, ?) на реальные значения
// Делает запросы читаемыми для логов и отладки
func Pretty(query string, placeholder string, args ...any) string {
	// Проходим по всем аргументам, переданным в запрос
	for i, param := range args {
		var value string

		// Определяем тип параметра и форматируем его соответствующим образом
		switch v := param.(type) {
		case string:
			// Для строк добавляем кавычки: "John" → ""John""
			value = fmt.Sprintf("%q", v)
		case []byte:
			// Для байтовых массивов (например, binary data) конвертируем в строку
			value = fmt.Sprintf("%q", string(v))
		default:
			// Для всех остальных типов (int, float, bool и т.д.) используем %v
			value = fmt.Sprintf("%v", v)
		}

		// Заменяем плейсхолдер (например, $1) на отформатированное значение
		// i+1 потому что индексы в SQL начинаются с 1, а не с 0
		query = strings.Replace(
			query, // исходная строка
			fmt.Sprintf("%s%s", placeholder, strconv.Itoa(i+1)), // что ищем: "$1"
			value, // на что заменяем: "John"
			-1,    // -1 = заменяем все вхождения
		)
	}

	// Удаляем лишние пробелы и переносы строк для компактности
	query = strings.ReplaceAll(query, "\t", "")  // Удаляем табуляции
	query = strings.ReplaceAll(query, "\n", " ") // Заменяем переносы строк на пробелы

	// Убираем лишние пробелы в начале и конце строки
	return strings.TrimSpace(query)
}

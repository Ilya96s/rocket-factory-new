package part

import (
	"context"
	"fmt"

	errs "github.com/Ilya96s/rocket-factory-new/inventory/internal/errors"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/converter"
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/record"
)

// List возвращает детали, соответствующие переданному фильтру.
//
// Поддерживаются три режима:
//
//  1. Если переданы UUID — возвращаются детали с этими UUID.
//     Фильтр по типу в этом случае игнорируется.
//
//  2. Если UUID не переданы, но указан тип — возвращаются детали этого типа.
//
//  3. Если UUID и тип не указаны — возвращаются все детали.
//
// При запросе по UUID порядок результата совпадает с порядком UUID в фильтре.
func (r *repository) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	// Общая часть SELECT для всех вариантов фильтрации.
	//
	// Явно перечислять колонки лучше, чем использовать SELECT *:
	// изменение структуры таблицы не сломает порядок аргументов в Scan.
	const selectColumns = `
		SELECT
			uuid,
			name,
			description,
			part_type,
			price,
			stock_quantity,
			created_at
		FROM parts`

	// По умолчанию выбираем все детали и сортируем их по названию.
	query := selectColumns + `ORDER BY name`

	// args содержит значения для SQL-плейсхолдеров:
	// $1, $2 и так далее.
	//
	// Если фильтров нет, слайс останется пустым.
	var args []any

	switch {
	case len(filter.UUIDs) > 0:
		// Если переданы UUID, фильтрация по типу игнорируется.
		//
		// ANY позволяет передать весь слайс UUID одним параметром:
		//
		//     uuid = ANY($1::uuid[])
		//
		// array_position сохраняет порядок, в котором UUID были
		// переданы клиентом. Без него PostgreSQL не гарантирует
		// порядок возвращаемых строк.
		query = selectColumns + `
			WHERE uuid = ANY($1::uuid[])
			ORDER BY array_position($1::uuid[], uuid)`

		args = append(args, filter.UUIDs)

	case filter.PartType != model.PartTypeUnspecified:
		// Если UUID отсутствуют, но указан тип детали,
		// фильтруем записи по колонке part_type.
		query = selectColumns + `
			WHERE part_type = $1
			ORDER BY name`

		args = append(args, string(filter.PartType))
	}

	// Выполняем запрос непосредственно через pgxpool.
	//
	// query содержит SQL с плейсхолдерами,
	// args содержит значения для этих плейсхолдеров.
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"получить список деталей: %w",
			err,
		)
	}

	// Обязательно закрываем rows после завершения чтения.
	// Это возвращает соединение обратно в пул.
	defer rows.Close()

	// Заранее выделяем пустой слайс результата.
	parts := make([]model.Part, 0)

	for rows.Next() {
		// Record-модель представляет строку таблицы PostgreSQL.
		var partRecord record.Part

		// Scan переносит значения текущей строки результата
		// в поля record-модели.
		//
		// Порядок аргументов должен точно совпадать с порядком
		// колонок в SELECT.
		if err := rows.Scan(
			&partRecord.UUID,
			&partRecord.Name,
			&partRecord.Description,
			&partRecord.PartType,
			&partRecord.Price,
			&partRecord.StockQuantity,
			&partRecord.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"прочитать деталь из результата запроса: %w",
				err,
			)
		}

		// Преобразуем модель хранилища в доменную модель.
		part := converter.FromRecordToModel(partRecord)

		parts = append(parts, part)
	}

	// rows.Next прекращает итерацию как при достижении конца результата,
	// так и при возникновении ошибки.
	//
	// Поэтому после цикла нужно отдельно проверить rows.Err.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"прочитать строки списка деталей: %w",
			err,
		)
	}

	// Если выполнялся поиск по UUID, PostgreSQL не вернёт ошибку
	// для отсутствующего UUID — он просто вернёт меньше строк.
	//
	// Например, запросили три UUID, но в таблице существуют только два.
	// Для контракта InventoryService это означает ErrPartNotFound.
	if len(filter.UUIDs) > 0 &&
		len(parts) != len(filter.UUIDs) {
		return nil, errs.ErrPartNotFound
	}

	return parts, nil
}

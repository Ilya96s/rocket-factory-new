package errs

import "errors"

var (
	// Ошибки деталей
	ErrPartNotFound = errors.New("деталь не найдена")

	// Ошибки валидации
	ErrInvalidUUID = errors.New("неверный формат UUID")
)

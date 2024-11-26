package service

import (
	"l0/internal/model"
	"l0/internal/repository"
)

// принимает на вход указатель на объект model.Order и
// пытается сохранить этот объект в базе данных с помощью библиотеки GORM
// Если операция завершается с ошибкой, функция возвращает эту ошибку.
func CreateOrder(order *model.Order) error {
	return repository.DB.Create(order).Error
}

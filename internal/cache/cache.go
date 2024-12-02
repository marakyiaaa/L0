package cache

import (
	"fmt"
	"l0/internal/model"
	"sync"
)

type OrderCache interface {
	CreateOrder(order model.Order)
	GetOrder(id string) (model.Order, error)
}

type Cache struct {
	msg sync.Map
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) IsEmpty() bool {
	var isEmpty bool
	c.msg.Range(func(key, value interface{}) bool {
		isEmpty = false
		return false
	})
	return isEmpty
}

// добавление или обновления данных в кэше
func (c *Cache) CreateOrder(order model.Order) {
	c.msg.Store(order.Order_uid, order)
}

// получение данных по ключу
func (c *Cache) GetOrder(id string) (model.Order, error) {
	value, ok := c.msg.Load(id)
	if !ok {
		return model.Order{}, fmt.Errorf("заказс id %s не найден в кэше", id)
	}
	return value.(model.Order), nil
}

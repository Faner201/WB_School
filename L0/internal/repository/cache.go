package repository

import (
	"L0/internal/entity"
	"sync"
)

type Cache struct {
	mu   sync.Mutex
	data map[string]entity.Order
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]entity.Order),
	}
}

func (c *Cache) SetOrder(order entity.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func (c *Cache) GetOrderByUID(uid string) (entity.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	order, ok := c.data[uid]
	return order, ok
}

func (c *Cache) CleanCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]entity.Order)
}

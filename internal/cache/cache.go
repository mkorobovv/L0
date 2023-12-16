package cache

import (
	"fmt"
	"sync"

	"github.com/mkorobovv/L0/internal/database"
	"github.com/mkorobovv/L0/internal/models"
)

type Cache struct {
	mu   sync.RWMutex
	db   *database.DB
	data map[string]*models.OrderJSON
}

func NewCache(db *database.DB) *Cache {
	return &Cache{
		data: make(map[string]*models.OrderJSON),
		db:   db,
		mu:   sync.RWMutex{},
	}
}

func (c *Cache) GetOrderInfo(uid string) *models.OrderJSON {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order := c.data[uid]
	return order
}

func (c *Cache) SetOrderInfo(order models.OrderJSON) {
	err := c.db.SetOrder(order)
	if err != nil {
		fmt.Printf("Cannot insert order: %v\n", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.Order_uid] = &order
	fmt.Printf("Cache written: %s\n", order.Order_uid)
}

func (c *Cache) GetAllOrders() map[string]*models.OrderJSON {
	return c.data
}

func (c *Cache) Preload() {
	ord, err := c.db.GetAllOrders()
	if err != nil {
		fmt.Printf("Error at DB: %v\n", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, or := range ord {
		c.data[or.Order_uid] = &or
	}
}

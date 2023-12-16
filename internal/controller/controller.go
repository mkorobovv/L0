package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mkorobovv/L0/internal/cache"
)

type CacheHandler struct {
	c *cache.Cache
}

func NewOrderController(c *cache.Cache) *CacheHandler {
	return &CacheHandler{
		c: c,
	}
}

func (c *CacheHandler) GetOrderController(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	order := c.c.GetOrderInfo(ctx.Param("order"))
	ord, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		fmt.Printf("Error at marshaling respond: %v", err)
	}
	return ctx.JSONBlob(http.StatusOK, ord)
}

func (c *CacheHandler) GetAllOrders(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	order := c.c.GetAllOrders()

	ord, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling respond: %v", err)
	}

	return ctx.JSONBlob(http.StatusOK, ord)
}

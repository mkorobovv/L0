package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mkorobovv/L0/internal/cache"
)

func OrderHandler(ctx echo.Context, oc *cache.Cache) error {
	ctx.Response().Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(oc.GetOrderInfo(ctx.Param("uid")), "", "\t")

	if err != nil {
		fmt.Printf("Error at marshaling: %v", err)
		return err
	}

	return ctx.JSONBlob(http.StatusOK, response)
}

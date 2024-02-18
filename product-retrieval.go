package main

import "github.com/labstack/echo/v4"

func (h *Handler) getProducts(c echo.Context) error {
	productType := new(struct {
		prodcutTypeName string `json:"productTypeName`
	})
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) deleteProduct(c echo.Context) error {

	productTypeName := c.Param("productTypeName")
	productNumberParam := c.Param("productNumber")

	var productTypeId int
	productNumber, err := strconv.Atoi(productNumberParam)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Bad Integer Conversion"})
	}

	fmt.Println(productTypeName, productNumber)

	err = h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", productTypeName).Scan(&productTypeId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Product Type doesnt Exist"})
	}

	_, err = h.DB.Exec("SELECT COUNT(*) FROM products WHERE producttypeid=? and productnumber=?", productTypeId, productNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Product Number Does Not Exist"})
	}

	_, err = h.DB.Exec("DELETE FROM product_attribute_values WHERE producttypeid=? and productnumber=?", productTypeId, productNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Deleting product Attribute Values"})
	}

	_, err = h.DB.Exec("DELETE FROM products WHERE producttypeid=? and productnumber=?", productTypeId, productNumber)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Deleting Prodcut"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product Deleted"})
}

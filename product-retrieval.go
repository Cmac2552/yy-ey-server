package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) getProductAttributeNames(c echo.Context) error {
	productType := new(productTypeBody)
	c.Bind(productType)

	result, err := h.DB.Query("SELECT product_attribute.attributename FROM product_type JOIN producttypeAttributes ON product_type.id = producttypeAttributes.producttypeid JOIN product_attribute ON product_attribute.id = producttypeAttributes.productattributeid WHERE product_type.productname = ?",
		productType.ProductTypeName)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database Error"})
	}

	attributeNames := make([]string, 0)

	for result.Next() {
		var attributeName string
		err = result.Scan(&attributeName)
		defer result.Close()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
		}
		attributeNames = append(attributeNames, attributeName)

	}
	return c.JSON(http.StatusOK, echo.Map{"productAttributeName": attributeNames})
}

func (h *Handler) getProducts(c echo.Context) error {
	productType := new(productTypeBody)
	c.Bind(productType)

	result, err := h.DB.Query("SELECT product_attribute.attributename,  products.attributevalue FROM product_type JOIN products ON product_type.id = products.producttypeid JOIN product_attribute ON product_attribute.id = products.productattributeid WHERE product_type.productname = ?",
		productType.ProductTypeName)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database Error"})
	}

	attributeValues := make([]productValue, 0)

	for result.Next() {
		var attributeValue productValue
		err = result.Scan(&attributeValue)
		defer result.Close()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
		}
		attributeValues = append(attributeValues, attributeValue)

	}
	return c.JSON(http.StatusOK, echo.Map{"products": attributeValues})
}

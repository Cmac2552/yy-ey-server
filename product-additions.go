package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) addProductType(c echo.Context) error {
	p := new(productType)
	c.Bind(p)
	orgId := c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).OrgId
	var count int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM product_type WHERE productname = ? AND orgid = ?", p.ProductName, orgId).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database Error"})
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, echo.Map{"message": "Product Already Exists"})
	}

	_, err = h.DB.Exec("INSERT INTO product_type (productname, orgid) VALUES(?, ?)", p.ProductName, orgId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Creating Product"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product Created"})
}

func (h *Handler) addProductAttribute(c echo.Context) error {
	incomingProductAttribute := new(newProductAttribute)
	c.Bind(incomingProductAttribute)
	var existingCount int
	err := h.DB.QueryRow("SELECT COUNT(*) FROM product_attribute WHERE attributename = ?", incomingProductAttribute.AttributeName).Scan(&existingCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Fetcing Product Attributes"})
	}

	if existingCount < 1 {
		result, err := h.DB.Exec("INSERT INTO product_attribute (attributename) VALUES(?)", incomingProductAttribute.AttributeName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Creating Product Attribute"})
		}

		productAttributeId, err := result.LastInsertId()

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Getting Product Attribute ID"})
		}

		_, err = h.DB.Exec("INSERT INTO producttypeAttributes (producttypeid, productattributeid) VALUES((SELECT id FROM product_type WHERE productname = ?), ?)", incomingProductAttribute.ProductName, int(productAttributeId))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Creating Product Attribute Connection"})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Product Attribute Created"})
	} else {
		var productTypeId int
		var productAttributeId int
		var count int
		err := h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", incomingProductAttribute.ProductName).Scan(&productTypeId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error getting product id"})
		}

		err = h.DB.QueryRow("SELECT id FROM product_attribute WHERE attributename=?", incomingProductAttribute.AttributeName).Scan(&productAttributeId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error getting attribute id"})
		}

		err = h.DB.QueryRow("SELECT COUNT (*) FROM producttypeAttributes WHERE producttypeid=? AND productattributeid = ?", productTypeId, productAttributeId).Scan(&count)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error checking if this attribute exists on this product"})
		}

		if count < 1 {
			_, err = h.DB.Exec("INSERT INTO producttypeAttributes (producttypeid, productattributeid) VALUES(?, ?)", productTypeId, productAttributeId)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Creating Product Attribute Connection"})
			}

			return c.JSON(http.StatusOK, echo.Map{"message": "Product Attribute Created"})
		} else {
			return c.JSON(http.StatusConflict, echo.Map{"message": "Error Fechting Attributes"})

		}
	}

}

func (h *Handler) addProductAttributeValue(c echo.Context) error {
	incomingProductAttributeValue := new(productAttributeValue)
	c.Bind(incomingProductAttributeValue)
	var count int
	var productTypeId int
	var productAttributeId int

	err := h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", incomingProductAttributeValue.ProductName).Scan(&productTypeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error getting product id"})
	}

	err = h.DB.QueryRow("SELECT id FROM product_attribute WHERE attributename=?", incomingProductAttributeValue.AttributeName).Scan(&productAttributeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error getting attribute id"})
	}

	err = h.DB.QueryRow("SELECT COUNT(*) FROM products WHERE producttypeid=? AND productattributeid = ?", productTypeId, productAttributeId).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error checking if attribute has value"})
	}

	if count > 0 {
		return c.JSON(http.StatusConflict, echo.Map{"message": "Product Already Exists"})
	}

	_, err = h.DB.Exec("INSERT INTO products (producttypeid, productattributeid, attributevalue) VALUES(?, ?, ?)", productTypeId, productAttributeId, incomingProductAttributeValue.AttributeValue)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error adding Attribute value"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Attribute Value Addition Successful"})

}

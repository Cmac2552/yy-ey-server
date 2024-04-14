package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) addProductType(c echo.Context) error {
	productType := new(productTypeBody)
	c.Bind(productType)
	orgId := c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).OrgId
	httpCode, httpMessage := h.addProductTypeAction(orgId, productType.ProductTypeName)
	return c.JSON(httpCode, echo.Map{"message": httpMessage})
}

func (h *Handler) addProductAttribute(c echo.Context) error {
	productAttribute := new(struct {
		ProductName   string `json:"productTypeName"`
		AttributeName string `json:"attributeName"`
	})
	c.Bind(productAttribute)
	httpCode, httpMessage := h.addProductAttributeAction(productAttribute.ProductName, productAttribute.AttributeName)
	return c.JSON(httpCode, echo.Map{"message": httpMessage})

}

func (h *Handler) addProductAttributeValue(c echo.Context) error {
	productAttributeValue := new(struct {
		AttributeValue string `json:"attributeValue"`
		AttributeName  string `json:"attributeName"`
		ProductName    string `json:"productName"`
		ProductNumber  int    `json:"productNumber"`
	})
	c.Bind(productAttributeValue)

	httpCode, httpMessage := h.addProductAttributeValueAction(
		productAttributeValue.ProductNumber,
		productAttributeValue.ProductName,
		productAttributeValue.AttributeName,
		productAttributeValue.AttributeValue)

	return c.JSON(httpCode, echo.Map{"message": httpMessage})
}

func (h *Handler) addProduct(c echo.Context) error {
	productType := new(productTypeBody)
	c.Bind(productType)
	httpCode, httpMessage, _ := h.addProductAction(productType.ProductTypeName)
	return c.JSON(httpCode, echo.Map{"message": httpMessage})

}

func (h *Handler) addProductTypeAndAttributeTypes(c echo.Context) error {
	productTypeAndAttributes := new(productTypeAndAttributes)
	c.Bind(productTypeAndAttributes)
	orgId := c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).OrgId

	httpCode, httpMessage := h.addProductTypeAction(orgId, productTypeAndAttributes.ProductTypeName)
	if httpCode != 200 {
		return c.JSON(httpCode, echo.Map{"message": httpMessage})
	}

	for _, element := range productTypeAndAttributes.Attributes {
		httpCode, httpMessage := h.addProductAttributeAction(productTypeAndAttributes.ProductTypeName, element)
		if httpCode != 200 {
			return c.JSON(httpCode, echo.Map{"message": httpMessage})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product Type Created"})
}

func (h *Handler) addProductAndAttributeValues(c echo.Context) error {
	productTypeAndAttributeValues := new(productTypeAndAttributeValues)
	c.Bind(productTypeAndAttributeValues)

	httpCode, httpMessage, productNumber := h.addProductAction(productTypeAndAttributeValues.ProductTypeName)
	if httpCode != 200 {
		return c.JSON(httpCode, echo.Map{"message": httpMessage})
	}

	for attrName, attrValue := range productTypeAndAttributeValues.AttrValues {
		httpCode, httpMessage := h.addProductAttributeValueAction(productNumber, productTypeAndAttributeValues.ProductTypeName, attrName, attrValue)
		if httpCode != 200 {
			return c.JSON(httpCode, echo.Map{"message": httpMessage})
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Product and Attribute Values Created"})
}

func (h *Handler) updateProductAttributeValues(c echo.Context) error {
	productTypeAndNumberWithAttributeValues := new(productTypeAndNumberWithAttributeValues)
	c.Bind(productTypeAndNumberWithAttributeValues)

	for attrName, attrValue := range productTypeAndNumberWithAttributeValues.AttrValues {
		httpCode, httpMessage := h.updateProductAction(productTypeAndNumberWithAttributeValues.ProductNumber, productTypeAndNumberWithAttributeValues.ProductTypeName, attrName, attrValue)
		if httpCode != 200 {
			return c.JSON(httpCode, echo.Map{"message": httpMessage})
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Product Attribute Value Updated"})
}

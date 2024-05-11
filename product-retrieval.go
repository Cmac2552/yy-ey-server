package main

import (
	"cmp"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) getProductAttributeNames(c echo.Context) error {
	productTypeName := c.Param("productTypeName")
	h.lock.Lock()

	result, err := h.DB.Query("SELECT product_attribute.attributename FROM product_type JOIN producttypeAttributes ON product_type.id = producttypeAttributes.producttypeid JOIN product_attribute ON product_attribute.id = producttypeAttributes.productattributeid WHERE product_type.productname = ?",
		productTypeName)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database Error"})
	}
	h.lock.Unlock()

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

	productType := c.Param("productTypeName")

	h.lock.Lock()
	attributes, err := h.DB.Query("SELECT productnumber, attributevalue, product_attribute.attributename from product_attribute_values JOIN product_attribute ON product_attribute_values.productattributeid = product_attribute.id where producttypeid=(SELECT id FROM product_type WHERE productname=?) ORDER BY productnumber ", productType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
	}
	h.lock.Unlock()
	products := make(map[int]map[string]string)

	for attributes.Next() {
		var productNumber int
		var attributeValue string
		var attributeName string
		attributes.Scan(&productNumber, &attributeValue, &attributeName)
		if products[productNumber] == nil {
			products[productNumber] = make(map[string]string)
		}
		products[productNumber][attributeName] = attributeValue
	}

	flat := make([]map[string]string, 0)
	for key, value := range products {
		value["productNumber"] = strconv.Itoa(key)
		flat = append(flat, value)
	}

	slices.SortFunc(flat, func(i, j map[string]string) int {
		productNumber1, err := strconv.Atoi(i["productNumber"])
		if err != nil {
			log.Println(err)
		}
		productNumber2, err := strconv.Atoi(j["productNumber"])
		if err != nil {
			log.Println(err)
		}
		return cmp.Compare(productNumber1, productNumber2)
	})

	return c.JSON(http.StatusOK, echo.Map{"products": flat})
}

func (h *Handler) getFilteredProducts(c echo.Context) error {
	productType := c.Param("productTypeName")

	//create part of query to be inserted
	//should be post request

	h.lock.Lock()
	productNumbers, err := h.DB.Query("SELECT productnumber from product_attribute_values JOIN product_attribute ON product_attribute_values.productattributeid = product_attribute.id where producttypeid=(SELECT id FROM product_type WHERE productname=?)"+" GROUP BY productnumber ", productType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
	}
	h.lock.Unlock()
	productNumberQuery := ""
	for productNumbers.Next() {

		var productNumber int
		productNumbers.Scan(&productNumber)

		productNumberQuery = productNumberQuery + "OR productnumber=" + strconv.Itoa(productNumber)
	}

	if productNumberQuery != "" {
		productNumberQuery = productNumberQuery[3:]
	}

	attributes, err := h.DB.Query("SELECT productnumber, attributevalue, product_attribute.attributename from product_attribute_values JOIN product_attribute ON product_attribute_values.productattributeid = product_attribute.id where producttypeid=(SELECT id FROM product_type WHERE productname=?) AND("+productNumberQuery+") ORDER BY productnumber ", productType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
	}
	h.lock.Unlock()
	products := make(map[int]map[string]string)

	for attributes.Next() {
		var productNumber int
		var attributeValue string
		var attributeName string
		attributes.Scan(&productNumber, &attributeValue, &attributeName)
		if products[productNumber] == nil {
			products[productNumber] = make(map[string]string)
		}
		products[productNumber][attributeName] = attributeValue
	}

	flat := make([]map[string]string, 0)
	for key, value := range products {
		value["productNumber"] = strconv.Itoa(key)
		flat = append(flat, value)
	}

	slices.SortFunc(flat, func(i, j map[string]string) int {
		productNumber1, err := strconv.Atoi(i["productNumber"])
		if err != nil {
			log.Println(err)
		}
		productNumber2, err := strconv.Atoi(j["productNumber"])
		if err != nil {
			log.Println(err)
		}
		return cmp.Compare(productNumber1, productNumber2)
	})

	return c.JSON(http.StatusOK, echo.Map{"products": flat})
}

func (h *Handler) getProductNames(c echo.Context) error {

	orgId := c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).OrgId
	h.lock.Lock()
	result, err := h.DB.Query("SELECT productname FROM product_type WHERE orgid=?",
		orgId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database Error"})
	}
	h.lock.Unlock()

	productNames := make([]string, 0)

	for result.Next() {
		var productName string
		err = result.Scan(&productName)
		defer result.Close()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
		}
		productNames = append(productNames, productName)

	}
	return c.JSON(http.StatusOK, echo.Map{"productNames": productNames})
}

func (h *Handler) getProductFilters(c echo.Context) error {

	productType := c.Param("productTypeName")

	h.lock.Lock()
	attributes, err := h.DB.Query("SELECT attributevalue, product_attribute.attributename from product_attribute_values JOIN product_attribute ON product_attribute_values.productattributeid = product_attribute.id where producttypeid=(SELECT id FROM product_type WHERE productname=?) GROUP BY attributevalue", productType)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Error Reading Rows From DB"})
	}
	h.lock.Unlock()
	productFilters := make(map[string][]string)

	for attributes.Next() {
		var attributeValue string
		var attributeName string
		attributes.Scan(&attributeValue, &attributeName)
		if productFilters[attributeName] == nil {
			productFilters[attributeName] = make([]string, 1)
			productFilters[attributeName][0] = attributeValue
		} else {
			productFilters[attributeName] = append(productFilters[attributeName], attributeValue)
		}
	}

	return c.JSON(http.StatusOK, productFilters)
}

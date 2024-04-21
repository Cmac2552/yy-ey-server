package main

import (
	"log"
	"net/http"
)

func (h *Handler) addProductTypeAction(orgId int, productTypeName string) (int, string) {
	var count int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM product_type WHERE productname = ? AND orgid = ?", productTypeName, orgId).Scan(&count)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Database Error"
	}

	if count > 0 {
		return http.StatusConflict, "Product Already Exists"
	}

	_, err = h.DB.Exec("INSERT INTO product_type (productname, orgid) VALUES(?, ?)", productTypeName, orgId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error Creating Product"
	}
	return http.StatusOK, "Product Type Created"
}

func (h *Handler) addProductAttributeAction(productName string, attributeName string) (int, string) {
	var existingCount int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM product_attribute WHERE attributename = ?", attributeName).Scan(&existingCount)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error Fetcing Product Attributes"
	}

	if existingCount < 1 {
		result, err := h.DB.Exec("INSERT INTO product_attribute (attributename) VALUES(?)", attributeName)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error Creating Product Attribute"
		}

		productAttributeId, err := result.LastInsertId()

		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error Getting Product Attribute ID"
		}

		_, err = h.DB.Exec("INSERT INTO producttypeAttributes (producttypeid, productattributeid) VALUES((SELECT id FROM product_type WHERE productname = ?), ?)", productName, int(productAttributeId))
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error Creating Product Attribute Connection"
		}

		return http.StatusOK, "Product Attribute Created"
	} else {
		var productTypeId int
		var productAttributeId int
		var count int
		err := h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", productName).Scan(&productTypeId)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error getting product id"
		}

		err = h.DB.QueryRow("SELECT id FROM product_attribute WHERE attributename=?", attributeName).Scan(&productAttributeId)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error getting attribute id"
		}

		err = h.DB.QueryRow("SELECT COUNT (*) FROM producttypeAttributes WHERE producttypeid=? AND productattributeid = ?", productTypeId, productAttributeId).Scan(&count)
		if err != nil {
			log.Println(err)
			return http.StatusInternalServerError, "Error checking if this attribute exists on this product"
		}

		if count < 1 {
			_, err = h.DB.Exec("INSERT INTO producttypeAttributes (producttypeid, productattributeid) VALUES(?, ?)", productTypeId, productAttributeId)
			if err != nil {
				log.Println(err)
				return http.StatusInternalServerError, "Error Creating Product Attribute Connection"
			}

			return http.StatusOK, "Product Attribute Created"
		} else {
			return http.StatusConflict, "Error Fechting Attributes"

		}
	}
}

func (h *Handler) addProductAttributeValueAction(productNumber int, productName string, attributeName string, attributeValue string) (int, string) {
	var count int
	var productTypeId int
	var productAttributeId int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM products WHERE productnumber = ?", productNumber).Scan(&count)
	if count < 0 || err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "No Product With that number"
	}

	err = h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", productName).Scan(&productTypeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error getting product id"
	}

	err = h.DB.QueryRow("SELECT id FROM product_attribute WHERE attributename=?", attributeName).Scan(&productAttributeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error getting attribute id"
	}

	err = h.DB.QueryRow("SELECT COUNT(*) FROM product_attribute_values WHERE producttypeid=? AND productattributeid = ? AND productnumber = ? ", productTypeId, productAttributeId, productNumber).Scan(&count)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error checking if attribute has value"
	}

	if count > 0 {
		return http.StatusConflict, "Product Already Exists With Attribute"
	}

	_, err = h.DB.Exec("INSERT INTO product_attribute_values (productnumber, producttypeid, productattributeid, attributevalue) VALUES(?, ?, ?, ?)", productNumber, productTypeId, productAttributeId, attributeValue)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error adding Attribute value"
	}
	return http.StatusOK, "Attribute Value Addition Successful"
}

func (h *Handler) addProductAction(productTypeName string) (int, string, int) {
	var productTypeId int

	err := h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", productTypeName).Scan(&productTypeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error getting product id", 0
	}
	_, err = h.DB.Exec("INSERT INTO products (productnumber, producttypeid) VALUES((SELECT productnumber FROM products WHERE producttypeid = ? ORDER BY productnumber DESC LIMIT 1 )+1, ?)", productTypeId, productTypeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error Creating Product", 0
	}

	var productNumber int
	h.DB.QueryRow("SELECT productNumber FROM products WHERE producttypeid = ? ORDER BY productnumber DESC LIMIT 1", &productTypeId).Scan(&productNumber)

	return http.StatusOK, "Product Addition Successful", productNumber
}

func (h *Handler) updateProductAction(productNumber int, productName string, attributeName string, attributeValue string) (int, string) {
	var count int
	var productTypeId int
	var productAttributeId int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM products WHERE productnumber = ?", productNumber).Scan(&count)
	if count < 0 || err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "No Product With that number"
	}

	err = h.DB.QueryRow("SELECT id FROM product_type WHERE productname=?", productName).Scan(&productTypeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error getting product id"
	}

	err = h.DB.QueryRow("SELECT id FROM product_attribute WHERE attributename=?", attributeName).Scan(&productAttributeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error getting attribute id"
	}

	err = h.DB.QueryRow("SELECT COUNT(*) FROM product_attribute_values WHERE producttypeid=? AND productattributeid = ? AND productnumber = ? ", productTypeId, productAttributeId, productNumber).Scan(&count)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error checking if attribute has value"
	}

	if count < 1 {
		return h.addProductAttributeValueAction(productNumber, productName, attributeName, attributeValue)
	}

	_, err = h.DB.Exec("UPDATE product_attribute_values SET attributevalue=? WHERE producttypeid=? AND productnumber=? AND productattributeid=?", attributeValue, productTypeId, productNumber, productAttributeId)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "Error adding Attribute value"
	}
	return http.StatusOK, "Attribute Value Update Successful"
}

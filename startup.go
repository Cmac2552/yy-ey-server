package main

import "fmt"

func (h *Handler) databaseStartUp() {
	_, err := h.DB.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		fmt.Println("Database Startup step 1 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS organizations (id INTEGER PRIMARY KEY AUTOINCREMENT, orgname TEXT)`)
	if err != nil {
		fmt.Println("Database Startup step 2 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT, orgid INTEGER REFERENCES organizations(id))`)
	if err != nil {
		fmt.Println("Database Startup step 3 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS product_type (id INTEGER PRIMARY KEY AUTOINCREMENT, productname TEXT, orgid INTEGER REFERENCES organizations(id))`)
	if err != nil {
		fmt.Println("Database Startup step 4 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS product_attribute (id INTEGER PRIMARY KEY AUTOINCREMENT, attributename TEXT UNIQUE)`)
	if err != nil {
		fmt.Println("Database Startup step 5 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS producttypeAttributes (id INTEGER PRIMARY KEY AUTOINCREMENT, producttypeid INTEGER REFERENCES product_type(id), productattributeid INTEGER REFERENCES product_attribute(id))`)
	if err != nil {
		fmt.Println("Database Startup step 6 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY AUTOINCREMENT, productnumber INTEGER, producttypeid INTEGER REFERENCES product_type(id))`)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Database Startup step 7 Failed")
	}

	_, err = h.DB.Exec(`CREATE TABLE IF NOT EXISTS product_attribute_values (id INTEGER PRIMARY KEY AUTOINCREMENT, productnumber INTEGER REFERENCES products(id), producttypeid INTEGER REFERENCES product_type(id),  productattributeid INTEGER REFERENCES product_attribute(id), attributevalue TEXT)`)
	if err != nil {
		fmt.Println("Database Startup step 8 Failed")
	}

}

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"
)

type user struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Organization string `json:"organization"`
}

type Item struct {
	ProductType string `json:"productType"`
	Descriptors string `json:"descriptors"`
}

type Handler struct {
	DB *sql.DB
}

func (h *Handler) restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func (h *Handler) addItem(c echo.Context) error {
	item := new(Item)
	c.Bind(item)
	userInfo := c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims)
	fmt.Println(item.ProductType, item.Descriptors, userInfo.Organization)
	_, err := h.DB.Exec("INSERT INTO products (producttype, descriptors, organization) VALUES (?, ?, ?)", item.ProductType, item.Descriptors, userInfo.Organization)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "Item Insertion failed"})
	} else {
		return c.JSON(http.StatusOK, echo.Map{"message": "Item Added"})
	}
}

func main() {
	e := echo.New()
	congfigureSecrets()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	db, err := sql.Open("sqlite", "./DB1.db")
	if err != nil {
		fmt.Println("error")
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS organizations (id INTEGER PRIMARY KEY AUTOINCREMENT, orgname TEXT)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, password TEXT, orgid INTEGER, FOREIGN KEY(orgid)REFERENCES organizations(id))`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tablez (id INTEGER PRIMARY KEY AUTOINCREMENT, tablename TEXT, orgid INTEGER, FOREIGN KEY(orgid)REFERENCES organizations(id))`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS product (id INTEGER PRIMARY KEY AUTOINCREMENT, productname TEXT)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS product_attribute (id INTEGER PRIMARY KEY AUTOINCREMENT, attributename TEXT)`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tableProducts (id INTEGER PRIMARY KEY AUTOINCREMENT, tableid INTEGER, productid INTEGER, FOREIGN KEY(tableid)REFERENCES tablez(id), FOREIGN KEY(productid)REFERENCES product(id))`)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS productAttributes (id INTEGER PRIMARY KEY AUTOINCREMENT, productid INTEGER, productattributeid INTEGER, FOREIGN KEY(productid)REFERENCES product(id), FOREIGN KEY(productattributeid)REFERENCES product_attribute(id))`)
	fmt.Println(err)
	h := &Handler{DB: db}

	// Login route
	e.POST("/login", h.login)
	e.POST("/sign-up", h.signUp)

	// Restricted group
	r := e.Group("")
	r.Use(echojwt.WithConfig(config))
	r.GET("/restricted", h.restricted)

	i := r.Group("/inventory")
	// i.POST("/add-inventory", h.addInvetoryTable)
	i.POST("/add-item", h.addItem)
	// i.GET("/yaks")

	e.Logger.Fatal(e.Start(":1323"))
}

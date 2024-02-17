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
	Username string `json:"username"`
	Password string `json:"password"`
	OrgId    int    `json:"organization"`
}

type signUpInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OrgName  string `json:"organization"`
}

type productType struct {
	ProductName string `json:"productName"`
}

type productAttribute struct {
	AttributeName string `json:"attributeName"`
}

type newProductAttribute struct {
	ProductName   string `json:"productName"`
	AttributeName string `json:"attributeName"`
}

type productAndProductAttribute struct {
	ProductTypeId int
	AttributeId   int
}

type productAttributeValue struct {
	AttributeValue string `json:"attributeValue"`
	AttributeName  string `json:"attributeName"`
	ProductName    string `json:"productName"`
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
	h := &Handler{DB: db}
	h.databaseStartUp()

	// Login route
	e.POST("/login", h.login)
	e.POST("/sign-up", h.signUp)

	// Restricted group
	r := e.Group("")
	r.Use(echojwt.WithConfig(config))
	r.GET("/restricted", h.restricted)

	//inventory group
	i := r.Group("/inventory")
	i.POST("/add-product-type", h.addProductType)
	i.POST("/add-product-attribute", h.addProductAttribute)
	i.POST("/add-product-attribute-value", h.addProductAttributeValue)

	e.Logger.Fatal(e.Start(":1323"))
}

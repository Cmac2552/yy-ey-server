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
		fmt.Println(err)
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
	i.POST("/product", h.addProduct)
	i.POST("/product-and-attributes", h.addProductTypeAndAttributeTypes)
	i.POST("/product-and-attribute-values", h.addProductAndAttributeValues)
	i.PATCH("/product-and-attribute-values", h.updateProductAttributeValues)
	i.GET("/products-attribute-names/:productTypeName", h.getProductAttributeNames)
	i.GET("/products/:productTypeName", h.getProducts)

	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (h *Handler) signUp(c echo.Context) error {
	u := new(user)
	c.Bind(u)
	queryStatement := "SELECT username FROM users WHERE username = ?"
	var returnUsername string
	err := h.DB.QueryRow(queryStatement, u.Username).Scan(&returnUsername)

	switch {
	case err == sql.ErrNoRows:
		bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
		_, err = h.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", u.Username, string(bytes))
		if err != nil {
			fmt.Println(err)
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "User: " + u.Username + " successfully inserted"})
	case err != nil:
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "Database failed"})
	default:
		return c.JSON(http.StatusConflict, echo.Map{"message": "User already exists"})
	}

}

func (h *Handler) login(c echo.Context) error {
	u := new(user)
	c.Bind(u)

	queryStatement := "SELECT username, password FROM users WHERE username = ?"
	var returnUser user

	err := h.DB.QueryRow(queryStatement, u.Username).Scan(&returnUser.Username, &returnUser.Password)
	passwordMatch := bcrypt.CompareHashAndPassword([]byte(returnUser.Password), []byte(u.Password))

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User does not exists"})
	} else if u.Username != returnUser.Username || passwordMatch != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Username and password do not match"})
	}

	claims := &jwtCustomClaims{
		returnUser.Username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create my Own secrete
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

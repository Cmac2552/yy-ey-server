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
	Name         string `json:"name"`
	Admin        bool   `json:"admin"`
	Organization string `json:"organization"`
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
		res, err := h.DB.Exec("INSERT INTO organizations (orgname) VALUES (?)", u.Organization)
		if err != nil {
			fmt.Println(err)
		}

		lastInsertedId, err := res.LastInsertId()
		if err != nil {
			fmt.Println(err)
		}

		_, err = h.DB.Exec("INSERT INTO users (username, password, orgid) VALUES (?, ?, ?)", u.Username, string(bytes), lastInsertedId)
		if err != nil {
			fmt.Println(err)
		}

		claims := &jwtCustomClaims{
			u.Username,
			true,
			u.Organization,
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

	queryStatement := "SELECT username, password, orgid FROM users WHERE username = ?"
	var returnUser user
	var orgId int

	err := h.DB.QueryRow(queryStatement, u.Username).Scan(&returnUser.Username, &returnUser.Password, &orgId)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database error"})
	}
	passwordMatch := bcrypt.CompareHashAndPassword([]byte(returnUser.Password), []byte(u.Password))

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User does not exists"})
	} else if u.Username != returnUser.Username || passwordMatch != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Username and password do not match"})
	}

	queryStatement = "SELECT orgname FROM organizations WHERE id = ?"
	err = h.DB.QueryRow(queryStatement, orgId).Scan(&returnUser.Organization)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database error"})
	}

	//DUPLICATE CODE
	claims := &jwtCustomClaims{
		returnUser.Username,
		true,
		returnUser.Organization,
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

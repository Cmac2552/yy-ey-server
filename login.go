package main

import (
	"database/sql"
	"fmt"
	"log"
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
	OrgId int    `json:"organization"`
	jwt.RegisteredClaims
}

type signUpInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OrgName  string `json:"organization"`
}

func (h *Handler) signUp(c echo.Context) error {
	s := new(signUpInfo)
	c.Bind(s)
	queryStatement := "SELECT username FROM users WHERE username = ?"
	var returnUsername string
	err := h.DB.QueryRow(queryStatement, s.Username).Scan(&returnUsername)

	switch {
	case err == sql.ErrNoRows:
		bytes, err := bcrypt.GenerateFromPassword([]byte(s.Password), 14)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "Password Encryption failed"})
		}

		res, err := h.DB.Exec("INSERT INTO organizations (orgname) VALUES (?)", s.OrgName)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "Organization creation failed"})
		}

		lastInsertedId, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "Organization id creation failed"})
		}

		_, err = h.DB.Exec("INSERT INTO users (username, password, orgid) VALUES (?, ?, ?)", s.Username, string(bytes), lastInsertedId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, echo.Map{"messsage": "User Creation Failed"})
		}

		claims := &jwtCustomClaims{
			s.Username,
			true,
			int(lastInsertedId),
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
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database error"})
	}
	passwordMatch := bcrypt.CompareHashAndPassword([]byte(returnUser.Password), []byte(u.Password))

	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User does not exists"})
	} else if u.Username != returnUser.Username || passwordMatch != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Username and password do not match"})
	}

	queryStatement = "SELECT id FROM organizations WHERE id = ?"
	err = h.DB.QueryRow(queryStatement, orgId).Scan(&returnUser.OrgId)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Database error"})
	}

	//DUPLICATE CODE
	claims := &jwtCustomClaims{
		returnUser.Username,
		true,
		returnUser.OrgId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create my Own secret
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

package rest

import (
	"net/http"

	"github.com/dexconv/hasin/retrieve/db"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func userRegister(c echo.Context) error {
	usr := user{}
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "could not recieve user info",
		})
	}

	if len(usr.Username) < 5 || len(usr.Password) < 5 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "username and password must be more that 5 characters",
		})
	}

	passcrypted, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.MinCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "password could not be encrypted",
		})
	}

	if err = db.CheckUserNameExist(usr.Username); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := db.ApplyRegistration(usr.Username, passcrypted); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "apply registration failed",
		})
	}

	Token, err := CreateJwt(usr.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not create token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user registered",
		"token":   Token,
	})
}

func userLogin(c echo.Context) error {
	usr := user{}
	if err := c.Bind(&usr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "could not recieve user info",
		})
	}

	if len(usr.Username) < 5 || len(usr.Password) < 5 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "username and password must be more that 5 characters",
		})
	}

	up, err := db.RetrievePass(usr.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not retrieve password",
		})
	}

	err = bcrypt.CompareHashAndPassword(up, []byte(usr.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "wrong password",
		})
	}

	Token, err := CreateJwt(usr.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "could not create token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "user logged in",
		"token":   Token,
	})
}

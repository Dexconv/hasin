package db

import (
	"fmt"

	"gorm.io/gorm"
)

type user struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

func CheckUserNameExist(username string) (err error) {
	usr := user{}
	err = DB.Where("username = ?", username).Find(&usr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	if usr.Username == "" {
		return nil
	}
	return fmt.Errorf("user %s already exists", username)
}

func ApplyRegistration(usr string, pass []byte) (err error) {
	err = DB.Create(&user{
		Username: usr,
		Password: pass,
	}).Error
	return
}

func RetrievePass(username string) ([]byte, error) {
	user := user{}
	DB.Where("username = ?", username).Find(&user)
	if user.Username == "" {
		return nil, fmt.Errorf("error trying to retrive pass, user doesn't exist")
	}
	return user.Password, nil
}

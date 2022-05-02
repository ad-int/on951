package data_generator

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"on951/application"
	dbStructure "on951/database/structure"
)

func GenerateUser(name string, password string, cost int) error {

	log.Println("cost", cost)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return errors.New("unable to generate token")
	}

	user := dbStructure.User{
		Name:     name,
		Password: string(hash),
	}
	tx := application.GetApplication().
		GetDatabase().
		GetDB().
		Debug().
		FirstOrCreate(&user)
	if tx.RowsAffected < 1 {
		return errors.New("could not find or create a user")
	}
	return nil
}

package main

import (
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/jinzhu/gorm"
)

var db gorm.DB

// MakeDB fills the db variable with an actual database connection.
func MakeDB(username, password, database, host, dbtype string) error {
	var err error
	db, err = gorm.Open(dbtype, fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, database))

	if err != nil {
		return err
	}

	db.AutoMigrate(&Repository{}, &Build{})

	return nil
}

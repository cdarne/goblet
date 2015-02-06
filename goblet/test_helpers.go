package goblet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"
)

type dbConfig struct {
	User, Password, Protocol, Host, Database string
	Port                                     int
}

func LoadDBConfig() (*dbConfig, error) {
	data, err := ioutil.ReadFile("db/config.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read the DB config file: %s\n", err))
	}

	cnf := dbConfig{}
	err = json.Unmarshal(data, &cnf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not load the DB Config: %s\n", err))
	}
	return &cnf, nil
}

func GetTestDB() gorm.DB {
	cnf, err := LoadDBConfig()
	So(err, ShouldBeNil)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s", cnf.User, cnf.Password, cnf.Protocol, cnf.Host, cnf.Port, cnf.Database))
	So(err, ShouldBeNil)
	return db
}

func WithTransaction(fn func(db *gorm.DB)) func() {
	return func() {
		db := GetTestDB()
		tx := db.Begin()
		So(tx.Error, ShouldBeNil)

		Reset(func() {
			So(tx.Rollback().Error, ShouldBeNil)
		})

		fn(tx)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/cdarne/goblet/goblet"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cnf, err := goblet.LoadDBConfig("db/config.json")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=latin1", cnf.User, cnf.Password, cnf.Protocol, cnf.Host, cnf.Port, "mysql"))
	if err != nil {
		log.Fatalf("Could not connect to DB: %s\n", err)
	}
	defer db.Close()

	stmts := []string{fmt.Sprintf("DROP DATABASE %s", cnf.Database), fmt.Sprintf("CREATE DATABASE %s", cnf.Database)}
	err = execMulti(db, stmts)
	if err != nil {
		log.Fatalf("Error while executing query: %s\n", err)
	}

	testDb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=latin1", cnf.User, cnf.Password, cnf.Protocol, cnf.Host, cnf.Port, cnf.Database))
	if err != nil {
		log.Fatalf("Could not connect to test DB: %s\n", err)
	}
	defer testDb.Close()

	data, err := ioutil.ReadFile("db/schema.sql")
	if err != nil {
		log.Fatalf("Could not read the DB schema file: %s\n", err)
	}
	stmts = strings.Split(string(data), ";")
	err = execMulti(testDb, stmts)
	if err != nil {
		log.Fatalf("Error while loading the DB schema: %s\n", err)
	}
}

func execMulti(db *sql.DB, stmts []string) error {
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) > 0 {
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

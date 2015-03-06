package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/cdarne/goblet/goblet"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("goblet", "Tests framework utility")
	app.Command("db", "Database jobs", func(dbCmd *cli.Cmd) {
		dbCmd.Command("create", "Create test DB", dbCreate)
	})
	app.Run(os.Args)
}

func dbCreate(cmd *cli.Cmd) {
	dbConfig := cmd.StringOpt("c config", "db/config.json", "Path to the DB config file")
	dbSchema := cmd.StringOpt("s schema", "db/schema.sql", "Path to the DB schema file")

	cmd.Action = func() {
		cnf, err := goblet.LoadDBConfig(*dbConfig)
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

		data, err := ioutil.ReadFile(*dbSchema)
		if err != nil {
			log.Fatalf("Could not read the DB schema file: %s\n", err)
		}
		stmts = strings.Split(string(data), ";")
		err = execMulti(testDb, stmts)
		if err != nil {
			log.Fatalf("Error while loading the DB schema: %s\n", err)
		}
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

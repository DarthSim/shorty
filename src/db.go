package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/peterhellberg/env"
)

// DB tools ====================================================================

var db *sql.DB

func confirmDbReset() {
	msg := fmt.Sprintf("Do you really want to reset DB?")
	if !confirm(msg) {
		halt()
	}
}

func initDBSchema() {
	schema := `
		DROP TABLE IF EXISTS urls;

		CREATE SEQUENCE urls_id_seq;
		CREATE SEQUENCE urls_code_seq START 2704 MINVALUE 2704;

		CREATE TABLE IF NOT EXISTS urls (
			id int4 PRIMARY KEY NOT NULL DEFAULT nextval('urls_id_seq'),
			url varchar(255) NOT NULL,
			code varchar(32) NOT NULL,
			open_count int4 NOT NULL DEFAULT 0
		);

		CREATE UNIQUE INDEX urls_code_ind ON urls (code);

		ALTER SEQUENCE urls_id_seq OWNED BY urls.id;
		ALTER SEQUENCE urls_code_seq OWNED BY urls.code;`

	_, err := db.Exec(schema)
	checkErr(err)
}

func initDB(resetConfirmation bool) {
	var err error

	connectionString := env.String("DB_CONN", "dbname=shorty_dev sslmode=disable")

	db, err = sql.Open("postgres", connectionString)
	checkErr(err)

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if env.Bool("RESET_DB", false) {
		if resetConfirmation {
			confirmDbReset()
		}
		initDBSchema()
	}

	prepareDbQueries()
}

func closeDB() {
	db.Close()
}

// end of DB tools

// DB queries ==================================================================

var (
	dbqCodeSeq      *sql.Stmt
	dbqCreateUrl    *sql.Stmt
	dbqGetUrl       *sql.Stmt
	dbqHitRedirect  *sql.Stmt
	dbqGetOpenCount *sql.Stmt
)

func prepareDbQueries() (err error) {
	dbqCodeSeq, err = db.Prepare("SELECT nextval('urls_code_seq')")
	checkErr(err)

	dbqCreateUrl, err = db.Prepare("INSERT INTO urls (url, code) VALUES ($1, $2)")
	checkErr(err)

	dbqGetUrl, err = db.Prepare("SELECT url FROM urls WHERE code = $1 LIMIT 1")
	checkErr(err)

	dbqHitRedirect, err = db.Prepare("UPDATE urls SET open_count = open_count + 1 WHERE code = $1")
	checkErr(err)

	dbqGetOpenCount, err = db.Prepare("SELECT open_count FROM urls WHERE code = $1 LIMIT 1")
	checkErr(err)

	return
}

func createUrl(url string) (code string, err error) {
	var codeSeq int64

	if err = dbqCodeSeq.QueryRow().Scan(&codeSeq); err != nil {
		return
	}

	code = buildCode(codeSeq)
	_, err = dbqCreateUrl.Exec(url, code)

	return
}

func getUrl(code string) (url string, err error) {
	err = dbqGetUrl.QueryRow(code).Scan(&url)
	return
}

func hitRedirect(code string) (err error) {
	_, err = dbqHitRedirect.Exec(code)
	return
}

func getOpenCount(code string) (count int64, err error) {
	err = dbqGetOpenCount.QueryRow(code).Scan(&count)
	return
}

// end of DB queries

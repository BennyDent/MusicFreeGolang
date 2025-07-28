package main

import (
	"database/sql"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"musicfree.root/handlers"
)

type myBD sql.DB

func main() {

	var db, err = sql.Open("mysql", SqlInitialize().FormatDSN())
	http.Handle("/create/albumn", handlers.AlbumnCreate(db))
	http.ListenAndServe(":7190", nil)

}

func SqlInitialize() *mysql.Config {

	cfg := mysql.NewConfig()
	//os.Getenv("DBUSER")
	//os.Getenv("DBPASS")
	cfg.User = "root"
	cfg.Passwd = "saharok2342saharok"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "recordings"
	return cfg

}

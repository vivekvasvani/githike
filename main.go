package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vivekvasvani/githike/server"
)

var db *sql.DB

func main() {
	wait := make(chan struct{})
	getDB()
	server.NewServer(db)
	<-wait

}

func getDB() {
	var err error
	db, err = sql.Open("mysql", "root:hike@tcp(localhost:3306)/githike")
	if err != nil {
		fmt.Print(err.Error())
		panic("Not able to Connect To DataBase")
	}
	fmt.Println(db)
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Print("Error :", err)
	}
}

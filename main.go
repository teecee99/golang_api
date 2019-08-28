package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "golang_db"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8&parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")

	if r.Method == http.MethodGet {
		if nId == "0" {
			io.WriteString(w, "id =  0 ="+nId)
		} else {
			io.WriteString(w, "number = "+nId)
		}
	} else if r.Method == http.MethodPost {

		db := dbConn()
		selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
		if err != nil {
			panic(err.Error())
		}
		emp := Employee{}
		for selDB.Next() {
			var id int
			var name, city string
			err = selDB.Scan(&id, &name, &city)
			if err != nil {
				panic(err.Error())
			}
			emp.Id = id
			emp.Name = name
			emp.City = city
		}

		if emp.Id == 0 {
			io.WriteString(w, "ไม่พบข้อมูล")
		} else {
			json.NewEncoder(w).Encode(emp)
		}
		defer db.Close()
	} else {
		io.WriteString(w, "This is a "+r.Method+" request")
	}

}
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the HomePage!")
}

func test(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")

	if r.Method == http.MethodGet {
		io.WriteString(w, "This is a get request"+nId)
	} else if r.Method == http.MethodPost {
		io.WriteString(w, "This is a post request")
	} else {
		io.WriteString(w, "This is a "+r.Method+" request")
	}
}

func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/test", test)
	http.HandleFunc("/employee", getEmployee)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleRequest()
}

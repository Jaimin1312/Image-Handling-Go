package main

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Id          int
	Name        string
	Image       string
	Author      string
	Price       string
	Isbn        string
	Language    string
	Description string
}

var tmpl = template.Must(template.ParseGlob("template/*"))

func dbConn() *sql.DB {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "Book"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database is connected")
	}
	return db
}

func admin(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM book ORDER BY id ASC")
	if err != nil {
		panic(err.Error())
	}

	res := []Book{}
	for selDB.Next() {
		var id int
		var image []byte
		var name, author, price, isbn, language, discription string
		err = selDB.Scan(&id, &name, &image, &author, &price, &isbn, &language, &discription)
		sEnc := b64.StdEncoding.EncodeToString(image)
		if err != nil {
			panic(err.Error())
		}

		res = append(res, Book{id, name, sEnc, author, price, isbn, language, discription})
	}
	file, _ := json.MarshalIndent(res, "", " ")
	_ = ioutil.WriteFile("Book.json", file, 0644)

	fmt.Println(res)
	tmpl.ExecuteTemplate(w, "admin.html", res)
	defer db.Close()

}

func main() {
	fmt.Println("Server started at 9000")
	http.HandleFunc("/admin", admin)
	http.ListenAndServe(":9000", nil)
}

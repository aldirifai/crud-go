package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)


type Product struct {
    Id    int
    Name  string
    Price int
    Description string
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Product
}

func Connect() (db *sql.DB) {
	dbDriver := "mysql"
    dbUser := "root"
    dbPass := ""
    dbName := "go_crud"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func main() {
	log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
	http.HandleFunc("/new", New)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/edit", Edit)
    http.HandleFunc("/insert", Insert)
    http.HandleFunc("/update", Update)
    http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}
var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	selDB, err := db.Query("SELECT * FROM products ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	p := Product{}
	res := []Product{}
	for selDB.Next() {
		var id, price int
		var name, description string
		err = selDB.Scan(&id, &name, &price, &description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = id
		p.Name = name
		p.Price = price
		p.Description = description
		res = append(res, p)
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Index", res)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM products WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := Product{}
	for selDB.Next() {
		var id, price int
		var name, description string
		err = selDB.Scan(&id, &name, &price, &description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = id
		p.Name = name
		p.Price = price
		p.Description = description
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Show", p)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	if r.Method == "POST" {
		name := r.FormValue("name")
		price := r.FormValue("price")
		description := r.FormValue("description")
		insForm, err := db.Prepare("INSERT INTO products(name, price, description) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, price, description)
		log.Println("INSERT: Name: " + name + " | Price: " + price + " | Description: " + description)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM products WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	p := Product{}
	for selDB.Next() {
		var id, price int
		var name, description string
		err = selDB.Scan(&id, &name, &price, &description)
		if err != nil {
			panic(err.Error())
		}
		p.Id = id
		p.Name = name
		p.Price = price
		p.Description = description
	}
	defer db.Close()
	tmpl.ExecuteTemplate(w, "Edit", p)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	if r.Method == "POST" {
		name := r.FormValue("name")
		price := r.FormValue("price")
		description := r.FormValue("description")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE products SET name=?, price=?, description=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, price, description, id)
		log.Println("UPDATE: Name: " + name + " | Price: " + price + " | Description: " + description)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}


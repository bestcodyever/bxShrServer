package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://codyduskin:@localhost/bxShr?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

// Boxes struct
type Boxes struct {
	ID          int
	Name        string
	Email       string
	Description string
	Long        float32
	Lat         float32
}

func main() {
	http.HandleFunc("/boxes", boxIndex)
	http.HandleFunc("/boxes/show", boxShowByID)
	http.HandleFunc("/boxes/create", createBox)
	http.ListenAndServe(":8080", nil)
}

// Get all Boxes
func boxIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM boxes")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	bxs := make([]Boxes, 0)
	for rows.Next() {
		bx := Boxes{}
		err = rows.Scan(&bx.ID, &bx.Name, &bx.Email, &bx.Description, &bx.Long, &bx.Lat)
		if err != nil {
			fmt.Print(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bxs = append(bxs, bx)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	js, err := json.Marshal(bxs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// Get Boxes by ID
func boxShowByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM boxes WHERE id = $1", id)

	bx := Boxes{}
	err := row.Scan(&bx.ID, &bx.Name, &bx.Email, &bx.Description, &bx.Long, &bx.Lat)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(bx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// Add new boxes
func createBox(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bx := Boxes{}
	bx.Name = r.FormValue("name")
	bx.Email = r.FormValue("email")
	bx.Description = r.FormValue("description")
	p1 := r.FormValue("long")
	p2 := r.FormValue("lat")

	// validate form values
	if bx.Name == "" || bx.Email == "" || bx.Description == "" || p1 == "" || p2 == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert form values
	f64, err := strconv.ParseFloat(p1, 32)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}
	bx.Long = float32(f64)

	// convert form values
	f65, err := strconv.ParseFloat(p2, 32)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	bx.Lat = float32(f65)

	_, err = db.Exec("INSERT INTO boxes (name, email, description, long, lat) VALUES ($1, $2, $3, $4, $5)", bx.Name, bx.Email, bx.Description, bx.Long, bx.Lat)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

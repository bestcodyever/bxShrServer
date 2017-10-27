package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

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
	http.HandleFunc("/", boxIndex)
	http.ListenAndServe(":8080", nil)
}

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

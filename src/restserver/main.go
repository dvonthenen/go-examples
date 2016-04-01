package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type account struct {
	Id       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type accounts []account

var port int
var address string

func listusers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host="+address+" user=dev password=vmware dbname=demo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//SELECT
	rows, err := db.Query("SELECT id, username, name, email FROM account")
	if err != nil {
		panic(err)
	}

	var accts accounts

	for rows.Next() {
		var id int
		var username string
		var name string
		var email string
		err = rows.Scan(&id, &username, &name, &email)
		if err != nil {
			continue
		}

		fmt.Println("id:", id, " username:", username, " name:", name,
			" email:", email)

		accts = append(accts, account{id, username, name, email})
	}

	response, err := json.MarshalIndent(accts, "", "  ")
	if err != nil {
		panic(err) //not expecting error... just a short cut
	}

	fmt.Println("response:", string(response))
	fmt.Fprintf(w, string(response))
}

func adduser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	accts := accounts{}
	if err := json.Unmarshal(body, &accts); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("postgres", "host="+address+" user=dev password=vmware dbname=demo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//We dont use for _, acct := range *accts { because we want to update the IDs
	for i := 0; i < len(accts); i++ {
		fmt.Printf("Username: %s, Name: %s, Email: %s\n", accts[i].Username, accts[i].Name, accts[i].Email)

		//INSERT
		var userid int
		err = db.QueryRow("INSERT INTO account (username, name, email, endpoint) VALUES ($1, $2, $3, $4) RETURNING id",
			accts[i].Username, accts[i].Name, accts[i].Email, r.URL.String()).Scan(&userid)
		if err != nil {
			panic(err)
		}

		fmt.Println("INSERT:", userid)
		accts[i].Id = userid
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(accts); err != nil {
		panic(err)
	}
}

func deleteuser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(w, "Deleting ID: %d", id)

	db, err := sql.Open("postgres", "host="+address+" user=dev password=vmware dbname=demo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//DELETE
	var userid int
	err = db.QueryRow("DELETE FROM account WHERE id = $1 RETURNING id", id).Scan(&userid)
	if err != nil {
		panic(err)
	}

	fmt.Println("DELETE:", userid)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	//define flags
	flag.IntVar(&port, "port", 9000, "the port in which to bind to")
	flag.StringVar(&address, "address", "127.0.0.1", "the postgres server in which to bind to")
	//parse
	flag.Parse()

	mux := mux.NewRouter()
	mux.HandleFunc("/user", listusers).Methods("GET")
	mux.HandleFunc("/user", adduser).Methods("POST")
	mux.HandleFunc("/user/{id}", deleteuser).Methods("DELETE")
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + strconv.Itoa(port))
}

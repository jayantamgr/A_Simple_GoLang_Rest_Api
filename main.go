package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"github.com/rs/cors"
)

/* 
  A person object with first name and last name to show all persons.
*/
  
type person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"` 
}

/*
  Slice below stores all the person objects.
*/
var personSlice []*person

/*
  A Struct which takes a json object at the post endpoint to create a person in the database.
*/

type personCreate struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"` 
	Age       int `json:"age`
}

// Create a person

func createPerson(w http.ResponseWriter, r *http.Request) {
    /* 
	  Takes a Json object and decodes the body and creates a person and inserts into database
	*/

	w.Header().Set("Content-Type", "application/json")
	var person personCreate
	_ = json.NewDecoder(r.Body).Decode(&person)

	connStr := "postgres://postgres:*******@localhost/exampledb?sslmode=disable"
	database, error := sql.Open("postgres", connStr)
	if error != nil {
		log.Fatal(error)
	}
	createPersonStatement := `INSERT INTO person (firstname, lastname, age)
	VALUES ($1, $2, $3)`
	rows, err := database.Query(createPersonStatement, person.Firstname, person.Lastname, person.Age)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
}

// Delete a person

func deletePerson(w http.ResponseWriter, r *http.Request) { 
    /* 
	   Deletes a person from the database with the given firstname.
	*/
	 w.Header().Set("Content-Type", "application/json")
	 params := mux.Vars(r)

	 firstname := ""

	 for _, items := range personSlice {
		 if items.Firstname == params["firstname"] {
			//copy := json.NewDecoder(r.Body).Decode(params["firstname"])
			firstname = params["firstname"]
		 }
	 }
	 connStr := "postgres://postgres:*******@localhost/exampledb?sslmode=disable"
	 database, error := sql.Open("postgres", connStr)
	 if error != nil {
		 log.Fatal(error)
	 }
	 deletePersonStatement := `DELETE FROM person WHERE firstname = $1;`
	 rows, err := database.Query(deletePersonStatement, firstname)
	 if err != nil {
		 fmt.Println(err)
	 }
	 defer rows.Close()
}

// Get all person

func getPersonEndPoint(w http.ResponseWriter, req *http.Request) {
    /*
	   Gets all the person from the person table by encoding to json
	*/
    
	connStr := "postgres://postgres:*******@localhost/exampledb?sslmode=disable"
	database, error := sql.Open("postgres", connStr)
	if error != nil {
		log.Fatal(error)
	}

	rows, err := database.Query("SELECT firstname, lastname FROM person ")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		c := new(person)
		if err := rows.Scan(&c.Firstname, &c.Lastname); err != nil {
			fmt.Println(err)
		}
		personSlice = append(personSlice, c)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	} 
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(personSlice); err != nil {
		fmt.Println(err)
	}
}


/*
   MAIN PROGRAM
      Handling all the routes with Mux Router
*/
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/person", getPersonEndPoint).Methods("GET")
	router.HandleFunc("/create", createPerson).Methods("POST")
	router.HandleFunc("/delete/{firstname}", deletePerson).Methods("DELETE")
	handler := cors.Default().Handler(router)
    http.ListenAndServe(":5000", handler)
}

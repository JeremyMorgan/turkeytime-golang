package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

// flags for db access (do we need this? )
var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "[PASSWORD]", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "[SERVER]", "the database server")
var user = flag.String("user", "[USER]", "the database user")

//FoodItem is an item of food
type FoodItem struct {
	ID          int
	Person      string
	Description string
}

//FoodItems is a list of items
type FoodItems []*FoodItem

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))
}

//Index is our main function to be called by root
func Index(w http.ResponseWriter, r *http.Request) {
	flag.Parse() // parse the command line args

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)

	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM [thanksgiving].[dbo].[FoodItems]")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	for rows.Next() {
		var id int
		var Description string
		var Person string

		if err := rows.Scan(&id, &Description, &Person); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("|%d|%s|%s|\n", id, Person, Description)

		//prepare a list of types to send
		FoodItems := FoodItems{
			{ID: id, Person: Person, Description: Description},
		}

		if err := json.NewEncoder(w).Encode(FoodItems); err != nil {
			panic(err)
		}
	}

	// header stuff

	//fmt.Fprintln(w, FoodItems)

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

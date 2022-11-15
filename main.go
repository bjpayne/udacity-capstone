package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

const database string = "app.db"

var db, _ = sql.Open("sqlite3", database)

type Customer struct {
	Id        uint16
	FirstName string
	LastName  string
	Email     string
	Role      string
	Phone     string
	Street    string
	City      string
	State     string
	Zip       string
	Contacted int
	CreatedAt string
}

func index(response http.ResponseWriter, request *http.Request) {
	rows, _ := db.Query("SELECT * FROM customers")

	var customers []Customer

	for rows.Next() {
		customer := Customer{}

		err := rows.Scan(
			&customer.Id,
			&customer.FirstName,
			&customer.LastName,
			&customer.Email,
			&customer.Phone,
			&customer.Role,
			&customer.Street,
			&customer.City,
			&customer.State,
			&customer.Zip,
			&customer.Contacted,
			&customer.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
		}

		customers = append(customers, customer)
	}

	encodedCustomers, err := json.Marshal(customers)

	if err != nil {
		log.Println(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(encodedCustomers)
}

func show(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	row := db.QueryRow("SELECT * FROM customers where id = ?", vars["id"])

	customer := Customer{}

	err := row.Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.Phone,
		&customer.Role,
		&customer.Street,
		&customer.City,
		&customer.State,
		&customer.Zip,
		&customer.Contacted,
		&customer.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	encodedCustomer, err := json.Marshal(customer)

	if err != nil {
		log.Println(err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(encodedCustomer)
}

func store(response http.ResponseWriter, request *http.Request) {
	contacted := request.FormValue("contacted")

	_contacted := 0

	if contacted != "" {
		_contacted = 1
	}

	customer := Customer{
		FirstName: request.FormValue("first_name"),
		LastName:  request.FormValue("last_name"),
		Email:     request.FormValue("email"),
		Role:      request.FormValue("role"),
		Phone:     request.FormValue("phone"),
		Street:    request.FormValue("street"),
		City:      request.FormValue("city"),
		State:     request.FormValue("state"),
		Zip:       request.FormValue("zip"),
		Contacted: _contacted,
	}

	insertResult, storeCustomerError := db.Exec(
		"INSERT INTO customers VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		customer.Role,
		customer.Phone,
		customer.Street,
		customer.City,
		customer.State,
		customer.Zip,
		customer.Contacted,
		time.Now(),
	)

	if storeCustomerError != nil {
		log.Fatal(storeCustomerError)
	}

	newCustomerId, _ := insertResult.LastInsertId()

	row := db.QueryRow("SELECT * FROM customers WHERE id = ?", newCustomerId)

	scanError := row.Scan(
		&customer.Id,
		&customer.FirstName,
		&customer.LastName,
		&customer.Email,
		&customer.Phone,
		&customer.Role,
		&customer.Street,
		&customer.City,
		&customer.State,
		&customer.Zip,
		&customer.Contacted,
		&customer.CreatedAt,
	)

	if scanError != nil {
		log.Fatal(scanError)
	}

	encodedCustomer, encodeCustomerError := json.Marshal(customer)

	if encodeCustomerError != nil {
		log.Println(encodeCustomerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(encodedCustomer)
}

func update(response http.ResponseWriter, request *http.Request) {

}

func remove(response http.ResponseWriter, request *http.Request) {

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/customers/{id}", show).Methods("GET")
	router.HandleFunc("/", store).Methods("POST")
	router.HandleFunc("/customers/{id}", update).Methods("PUT")
	router.HandleFunc("/", remove).Methods("DELETE")

	fmt.Println("Server is starting on port 3000...")
	err := http.ListenAndServe(":3000", router)
	log.Fatal(err)
}

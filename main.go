package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"time"
)

const database string = "app.db"

var db, _ = sql.Open("sqlite3", database)

type Customer struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Contacted bool   `json:"contacted"`
	CreatedAt string `json:"created_at"`
}

func home(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "static/index.html")
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
			fatal(response, err)
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
	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	customerId, _ := strconv.Atoi(vars["id"])

	customer := fetchCustomer(customerId)

	if customer.Id != customerId {
		response.WriteHeader(404)
	}

	encodedCustomer, err := json.Marshal(customer)

	if err != nil {
		log.Println(err)
	}

	response.Write(encodedCustomer)
}

func store(response http.ResponseWriter, request *http.Request) {
	input := Customer{}

	inputDecodeError := json.NewDecoder(request.Body).Decode(&input)

	if inputDecodeError != nil {
		response.WriteHeader(500)
		log.Fatal(inputDecodeError)
	}

	insertResult, storeCustomerError := db.Exec(
		"INSERT INTO customers VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		nil,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Role,
		input.Phone,
		input.Street,
		input.City,
		input.State,
		input.Zip,
		input.Contacted,
		time.Now(),
	)

	if storeCustomerError != nil {
		log.Fatal(storeCustomerError)
	}

	newCustomerId, _ := insertResult.LastInsertId()

	customer := fetchCustomer(int(newCustomerId))

	encodedCustomer, encodeCustomerError := json.Marshal(customer)

	if encodeCustomerError != nil {
		log.Println(encodeCustomerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(encodedCustomer)
}

func update(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	input := Customer{}

	inputDecodeError := json.NewDecoder(request.Body).Decode(&input)

	vars := mux.Vars(request)

	customerId, _ := strconv.Atoi(vars["id"])

	customer := fetchCustomer(customerId)

	if customerId != customer.Id {
		response.WriteHeader(404)
	}

	input.Id = customerId

	if inputDecodeError != nil {
		fatal(response, inputDecodeError)
	}

	query := `
		UPDATE customers 
		SET 
		    first_name = ?,
		    last_name = ?, 
		    email = ?,
		    role = ?,
		    phone = ?,
		    street = ?,
		    city = ?,
		    state = ?,
		    zip = ?,
		    contacted = ?
		WHERE id = ?
	`

	_, updateCustomerError := db.Exec(
		query,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Role,
		input.Phone,
		input.Street,
		input.City,
		input.State,
		input.Zip,
		input.Contacted,
		customerId,
	)

	if updateCustomerError != nil {
		fatal(response, updateCustomerError)
	}

	customer = fetchCustomer(customerId)

	encodedCustomer, encodeCustomerError := json.Marshal(customer)

	if encodeCustomerError != nil {
		fatal(response, encodeCustomerError)
	}

	response.Write(encodedCustomer)
}

func remove(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(request)

	customerId, _ := strconv.Atoi(vars["id"])

	customer := fetchCustomer(customerId)

	if customerId != customer.Id {
		response.WriteHeader(404)
	}

	query := "DELETE FROM customers WHERE id = ?"

	_, err := db.Exec(query, customerId)

	if err != nil {
		fatal(response, err)
	}

	message := map[string]string{"message": fmt.Sprintf("%s %s (%s) deleted", customer.FirstName, customer.LastName, customer.Email)}

	encodedMessage, _ := json.Marshal(message)

	response.Write(encodedMessage)
}

func fetchCustomer(customerId int) Customer {
	row := db.QueryRow("SELECT * FROM customers WHERE id = ?", customerId)

	customer := Customer{}

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
		if scanError.Error() == "sql: no rows in result set" {
			return customer
		}

		log.Fatal(scanError)
	}

	return customer
}

func fatal(response http.ResponseWriter, message any) {
	response.WriteHeader(500)
	log.Fatal(message)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/customers", index).Methods("GET")
	router.HandleFunc("/customers/{id}", show).Methods("GET")
	router.HandleFunc("/customers", store).Methods("POST")
	router.HandleFunc("/customers/{id}", update).Methods("PUT")
	router.HandleFunc("/customers/{id}", remove).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Server is starting on port 3000...")
	err := http.ListenAndServe(":3000", router)
	log.Fatal(err)
}

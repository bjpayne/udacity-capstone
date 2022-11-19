package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/customers", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(index)
	handler.ServeHTTP(rr, req)

	// Checks for 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getCustomers returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

func TestShowHandler(t *testing.T) {
	router := mux.NewRouter()

	req, err := http.NewRequest("GET", "/customers/4", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.HandleFunc("/customers/{id}", show).Methods("GET")
	router.ServeHTTP(rr, req)

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf(
			"Content-Type does not match: got %v want %v",
			ctype,
			"application/json",
		)
	}

	// Checks for 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"show returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}
}

var storedCustomer Customer

func TestStoreHandler(t *testing.T) {
	requestBody := strings.NewReader(`
		{
			"first_name": "First111",
			"last_name": "Last",
			"email": "test@test.com",
			"role": "customer",
			"phone": "(111) 222-3344",
			"street": "1234 Test St.",
			"city": "Test",
			"state": "TE",
			"zip": "12345-1111",
			"contacted": false
		}
	`)

	req, err := http.NewRequest("POST", "/customers", requestBody)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(store)
	handler.ServeHTTP(rr, req)

	decodeResponseError := json.NewDecoder(rr.Body).Decode(&storedCustomer)

	if decodeResponseError != nil {
		return
	}

	// Checks for 200 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("addCustomer returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

func TestUpdateHandler(t *testing.T) {
	router := mux.NewRouter()

	customerJson, _ := json.Marshal(storedCustomer)

	requestBody := strings.NewReader(string(customerJson))

	req, err := http.NewRequest("PUT", "/customers/4", requestBody)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.HandleFunc("/customers/{id}", update).Methods("PUT")
	router.ServeHTTP(rr, req)

	// Checks for 201 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("update returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

func TestRemoveHandler(t *testing.T) {
	router := mux.NewRouter()

	req, err := http.NewRequest("DELETE", "/customers/"+strconv.Itoa(storedCustomer.Id), nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.HandleFunc("/customers/{id}", remove).Methods("DELETE")
	router.ServeHTTP(rr, req)

	// Checks for 404 status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("deleteCustomer returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

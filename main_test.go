package main

import (
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests happy path of submitting a well-formed GET / request
func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

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

// Tests happy path of showing a customer
func TestShowHandler(t *testing.T) {
	router := mux.NewRouter()

	req, err := http.NewRequest("GET", "/customers/2", nil)

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

// Tests happy path of submitting a well-formed POST /customers request
func TestStoreHandler(t *testing.T) {
	requestBody := strings.NewReader(`
		{
			"first_name": "First",
			"last_name": "Last",
			"email": "test@test.com",
			"role": "customer",
			"phone": "111-222-3344",
			"street": "1234 test street",
			"city": "City",
			"State": "ST",
			"Zip": "12345-1111",
			"Contacted": true,
		}
	`)

	req, err := http.NewRequest("POST", "/customers", requestBody)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(store)
	handler.ServeHTTP(rr, req)

	// Checks for 201 status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("addCustomer returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Checks for JSON response
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content-Type does not match: got %v want %v",
			ctype, "application/json")
	}
}

// Tests unhappy path of deleting a user that doesn't exist
func TestDeleteCustomerHandler(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/customers/2", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(remove)
	handler.ServeHTTP(rr, req)

	// Checks for 404 status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("deleteCustomer returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

// Tests unhappy path of getting a user that doesn't exist
func TestGetCustomerHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/customers/e7847fee-3a0e-455e-b151-519bdb9851c7", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(show)
	handler.ServeHTTP(rr, req)

	// Checks for 404 status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("getCustomer returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

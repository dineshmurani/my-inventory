package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occurred while initializing the database.")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY(id)
		);`

	_, err := a.DB.Exec(createTableQuery)

	if err != nil {
		log.Fatal(err)
	}
}

// test get products

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
	log.Println("ClearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}

}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chair", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, Got: %v", "chair", m["name"])
	}
	log.Printf("%T", m["quantity"])

	if m["quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, Got: %v", 1.0, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	// adding a row to the table
	addProduct("connector", 10, 10)

	// following GET call will fetch the above row from the table.
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
	// until here.

	// DELETE the products
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
	// DELETE the products

	// GET the product again
	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
	// GET the product again
}

func TestUpdateProduct(t *testing.T) {
	//clear table to start with
	clearTable()

	//add product as below:
	addProduct("connector", 10, 10)

	//get a response from the table about the above add product.
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	// store 'response' and 'unmarshal' it for later.
	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	//change value about the rows.
	var product = []byte(`{"name":"connector", "quantity":1, "price":10}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	// compare oldValue and newValue to make sure update is correct.
	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id: %v, Got: %v", newValue["id"], oldValue["id"])
	}

	// compare name now between oldValue and newValue
	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected name: %v, Got: %v", newValue["name"], oldValue["name"])
	}

	// compare price
	if oldValue["price"] != newValue["price"] {
		t.Errorf("Expected price: %v, Got: %v", newValue["price"], oldValue["price"])
	}

	// compare quantity
	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("Expected quantity: %v, Got: %v", newValue["quantity"], oldValue["quantity"])
	}
}

/*
1>
 go doc httptest Newrecorder
package httptest // import "net/http/httptest"

func NewRecorder() *ResponseRecorder
    NewRecorder returns an initialized ResponseRecorder.

2>
% go doc httptest ResponseRecorder
package httptest // import "net/http/httptest"

type ResponseRecorder struct {
        // Code is the HTTP response code set by WriteHeader.
ResponseRecorder is an implementation of http.ResponseWriter that records
    its mutations for later inspection in tests.

3>
% go doc github.com/gorilla/mux ServeHTTP
package mux // import "github.com/gorilla/mux"

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request)
    ServeHTTP dispatches the handler registered in the matched route.

    When there is a match, the route variables can be retrieved calling
    mux.Vars(request).

4>

*/

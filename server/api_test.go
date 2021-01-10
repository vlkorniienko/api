package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var userData map[string]int
var userDataFail map[string]int
var api *API

// Mock data
func init() {
	api = New(9000)
	userData = make(map[string]int)
	userDataFail = make(map[string]int)
	userData["id"] = 42
	userDataFail["id"] = 7
}

//go test -v -run TestCorrectProductID
func TestCorrectProductID(t *testing.T) {
	j, _ := json.Marshal(userData)

	req, err := http.NewRequest("GET", "/", bytes.NewBuffer(j))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler := http.HandlerFunc(api.paymentHandler)
	mux.Handle("/", handler)
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// user should be created in db after this request
	expected := `{"alipay":`

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//go test -v -run TestInCorrectProductID
func TestInCorrectProductID(t *testing.T) {
	j, _ := json.Marshal(userDataFail)

	req, err := http.NewRequest("GET", "/", bytes.NewBuffer(j))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler := http.HandlerFunc(api.paymentHandler)
	mux.Handle("/", handler)
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != 400{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// user should be created in db after this request
	expected := `Product with this ID doesn't exist"`
	// fmt.Println([]byte(expected))
	// fmt.Println([]byte(rr.Body.String()))

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//go test -v -run TestErrorFromPaymentServer
func TestErrorFromPaymentServer(t *testing.T) {
	j, _ := json.Marshal(userData)

	req, err := http.NewRequest("GET", "/", bytes.NewBuffer(j))
	if err != nil {
		t.Fatal(err)
	}
	api.PaymentMethods = []string{ // payment methods option
		"card",   // valid
		"alipay", // valid
		"p24", 
	}
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	handler := http.HandlerFunc(api.paymentHandler)
	mux.Handle("/", handler)
	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != 200{
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// user should be created in db after this request
	expected := `downloadAppLink":"https://applestore/download?id=1000"`

	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
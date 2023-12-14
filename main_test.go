package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestSignAnswers(t *testing.T) {
	// Create a request for the SignAnswers handler
	requestBody := []byte(`{"questions":["question1", "question2"],"answers":["answer1", "answer2"]}`)
	req, err := http.NewRequest("POST", "/sign/testuser", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the SignAnswers handler function
	handler := http.HandlerFunc(SignAnswers)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("SignAnswers returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	fmt.Println("Response Body:", rr.Body.String())

	// Parse the response body only if the status code is OK
	if rr.Code == http.StatusOK {
		// Parse the response body and check if it contains the expected fields
		var response Signature
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		// Example assertion
		if response.Timestamp.IsZero() {
			t.Errorf("SignAnswers returned unexpected timestamp: got %v want a non-zero timestamp", response.Timestamp)
		}
	}
}

func TestVerifySignature(t *testing.T) {
	// Create a request for the VerifySignature handler
	req, err := http.NewRequest("GET", "/verify/testuser/testsignature", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the VerifySignature handler function
	handler := http.HandlerFunc(VerifySignature)
	handler.ServeHTTP(rr, req)

	// Check the status code and response body
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("VerifySignature returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	fmt.Println("Response Body:", rr.Body.String())

	// Parse the response body only if the status code is OK
	if rr.Code == http.StatusOK {
		// Parse the response body and check if it contains the expected fields
		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}

		// Example assertion
		if status, ok := response["status"].(string); !ok || status != "OK" {
			t.Errorf("VerifySignature returned unexpected status: got %v want 'OK'", status)
		}
	}
}

func TestSignAnswers_InvalidJSON(t *testing.T) {

  reqBody := "invalid json"

  req := httptest.NewRequest("POST", "/sign/user", strings.NewReader(reqBody))

  rr := httptest.NewRecorder()

  handler := http.HandlerFunc(SignAnswers)
  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusBadRequest {
    t.Errorf("expected 400 status, got %v", status)
  }

}

func TestSignAnswers_NoAuthHeader(t *testing.T) {
  
  reqBody := `{"questions":[],"answers":[]}`

  req := httptest.NewRequest("POST", "/sign/user", strings.NewReader(reqBody))

  rr := httptest.NewRecorder()

  handler := http.HandlerFunc(SignAnswers)
  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnauthorized {
    t.Errorf("expected 401 status, got %v", status) 
  }

}

func TestSignAnswers_InvalidToken(t *testing.T) {

  reqBody := `{"questions":[],"answers":[]}`
  req := httptest.NewRequest("POST", "/sign/user", strings.NewReader(reqBody))

  req.Header.Set("Authorization", "InvalidToken")

  rr := httptest.NewRecorder()  

  handler := http.HandlerFunc(SignAnswers)
  handler.ServeHTTP(rr, req)

  if status := rr.Code; status != http.StatusUnauthorized {
    t.Errorf("expected 401 status, got %v", status)
  }

}
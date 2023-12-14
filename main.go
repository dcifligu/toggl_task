package main

import (
	"encoding/json" // for JSON encoding/decoding
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go" // JWT library
	"github.com/gorilla/mux"
)

type Signature struct {
	User      string    `json:"user"`
	Signature string    `json:"signature"` 
	Answers   []string  `json:"answers"`
	Timestamp time.Time `json:"timestamp"`
}

var signatures = make(map[string]Signature) // storage for signatures
var mu sync.Mutex // for mutual exclusion locking 

var jwtSecret = []byte("your-secret-key") // JWT signing key

// Handler to sign a set of question answers
func SignAnswers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		Questions []string `json:"questions"`
		Answers   []string `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	signature := Signature{
		User:      user,
		Signature: generateSignature(user, request.Questions, request.Answers),
		Answers:   request.Answers,
		Timestamp: time.Now(),
	}

	mu.Lock()
	signatures[user] = signature
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signature)
}
// Handler to verify a signature
func VerifySignature(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]
	signature := params["signature"]

	mu.Lock()
	storedSignature, ok := signatures[user]
	mu.Unlock()

	if !ok {
		http.Error(w, "Signature not found", http.StatusNotFound)
		return
	}

	if storedSignature.Signature == signature {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "OK",
			"user":      user,
			"answers":   storedSignature.Answers,
			"timestamp": storedSignature.Timestamp,
		})
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// Generate JWT-based signature
func generateSignature(user string, questions, answers []string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      user,
		"questions": questions,
		"answers":   answers,
		"timestamp": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// Main function to initialize router and start server
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/sign/{user}", SignAnswers).Methods("POST")
	router.HandleFunc("/verify/{user}/{signature}", VerifySignature).Methods("GET")

	port := 8080
	fmt.Printf("Server running on :%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

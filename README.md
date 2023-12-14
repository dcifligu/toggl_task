# Signature Service

This is a REST API service for signing and verifying data.

## Features

- Sign a set of questions and answers with a JWT-based signature
- Verify a signature matches the original signed data
- Protect endpoints with JWT authentication

## Endpoints

### POST /sign/:user

Signs the provided questions and answers and returns a signature object containing:

- user: username
- signature: JWT token
- answers: original answers
- timestamp: signature creation time

Requires valid JWT auth token.

### GET /verify/:user/:signature

Verifies if the signature matches the originally signed data.

Returns:

- status: OK or error
- user: username
- answers: original answers
- timestamp: signature creation time

## Running the service

`go run main.go`

This will start the server on port 8080.

## Generating tokens

Use the token.go file to generate valid JWT tokens for auth.

## Testing

Unit tests are provided in main_test.go using the net/http/httptest package.

Run tests:
`go test`

## TODO

- Persistent storage of signatures
- Input validation
- More tests

## Libraries

- [gorilla/mux](https://github.com/gorilla/mux) - router
- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) - JWTs

# CMPS2242 Lab 4 — Middleware, Dependency Injection, Custom Response Writer
Name: Vance-Petillo
Branch: lab04-development

## Run Server
go run main.go

## Test with curl
curl -i http://localhost:4000/v1/healthcheck
curl -i http://localhost:4000/v1/books
curl -i http://localhost:4000/v1/books/42
curl -i -X POST http://localhost:4000/v1/books
curl -i -X DELETE http://localhost:4000/v1/books/42
curl -i http://localhost:4000/v1/notfound

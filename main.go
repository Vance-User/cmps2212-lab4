package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// Custom Response Writer
//
// This struct wraps the standard http.ResponseWriter so we can intercept
// the WriteHeader call and capture the HTTP status code before it is sent.
////////////////////////////////////////////////////////////////////////////////

type responseWriter struct {
	http.ResponseWriter // embed the real writer
	statusCode          int // captured status code
}

// WriteHeader intercepts the status code before forwarding it.
// This allows middleware to log the real status after the handler runs.
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

////////////////////////////////////////////////////////////////////////////////
// Logging Middleware
//
// This middleware wraps the entire router.
// It captures request duration and HTTP status code for every request.
////////////////////////////////////////////////////////////////////////////////

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Record start time to measure request duration
		start := time.Now()

		// Wrap the real writer with our custom responseWriter.
		// Pre-set statusCode to 200 (OK) in case handler never calls WriteHeader.
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Pass wrapped writer to next handler in the chain
		next.ServeHTTP(rw, r)

		// Log method, path, status code, and duration
		log.Printf("%s %s %d %v",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

////////////////////////////////////////////////////////////////////////////////
// Application Struct (Dependency Injection)
//
// Instead of using global variables, shared dependencies (like logger)
// are grouped inside this struct and accessed via the receiver.
////////////////////////////////////////////////////////////////////////////////

type application struct {
	logger *slog.Logger
}

////////////////////////////////////////////////////////////////////////////////
// Handlers (Book Catalogue API - Version 1)
////////////////////////////////////////////////////////////////////////////////

// GET /v1/healthcheck
// Returns 200 OK confirming server is alive.
func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "status: available\n")

	app.logger.Info("healthcheck handler called")
}

// GET /v1/books
// Returns 200 OK listing books (placeholder for now).
func (app *application) listBooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "list of books (coming soon)\n")

	app.logger.Info("listBooks handler called")
}

// GET /v1/books/{id}
// Returns 200 OK and displays the requested book ID.
func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get book with id: %s\n", id)

	app.logger.Info("getBook handler called", "id", id)
}

// POST /v1/books
// Returns 201 Created.
func (app *application) createBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "book created (coming soon)\n")

	app.logger.Info("createBook handler called")
}

// DELETE /v1/books/{id}
// Returns 204 No Content.
// NOTE: 204 responses must NOT include a body.
func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	w.WriteHeader(http.StatusNoContent)

	app.logger.Info("deleteBook handler called", "id", id)
}

////////////////////////////////////////////////////////////////////////////////
// main() — Application Setup
////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create a structured logger writing to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Inject logger into application struct (dependency injection)
	app := &application{
		logger: logger,
	}

	// Create a new HTTP router 
	mux := http.NewServeMux()

	// Register all versioned routes
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheck)
	mux.HandleFunc("GET /v1/books", app.listBooks)
	mux.HandleFunc("GET /v1/books/{id}", app.getBook)
	mux.HandleFunc("POST /v1/books", app.createBook)
	mux.HandleFunc("DELETE /v1/books/{id}", app.deleteBook)

	// Log structured startup message
	logger.Info("starting server", "addr", ":4000")

	// Wrap entire router with logging middleware
	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	log.Fatal(err)
}
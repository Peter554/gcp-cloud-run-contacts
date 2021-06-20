package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	dsn := os.Getenv("DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := NewServer(db)
	log.Fatal(http.ListenAndServe(":"+port, server))
}

type Contact struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Server struct {
	mux *http.ServeMux
	db  *sql.DB
}

func NewServer(db *sql.DB) *Server {
	mux := http.NewServeMux()
	server := &Server{mux, db}

	mux.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.Write([]byte("hello!"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.getContacts(w, r)
		case http.MethodPost:
			server.postContacts(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) getContacts(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("select id, name, email from contacts limit 100")
	if err != nil {
		internalServerError(w, err)
		return
	}
	contacts := make([]Contact, 0)
	for rows.Next() {
		var contact Contact
		err := rows.Scan(&contact.ID, &contact.Name, &contact.Email)
		if err != nil {
			internalServerError(w, err)
			return
		}
		contacts = append(contacts, contact)
	}
	if err = rows.Err(); err != nil {
		internalServerError(w, err)
		return
	}
	w.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(contacts); err != nil {
		log.Println(err)
	}
}

func (s *Server) postContacts(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var contact Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		// Bad request?
		internalServerError(w, err)
		return
	}
	row := s.db.QueryRow("insert into contacts (name, email) values ($1, $2) returning id", contact.Name, contact.Email)
	if err := row.Scan(&contact.ID); err != nil {
		internalServerError(w, err)
		return
	}
	w.Header().Add("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(contact); err != nil {
		log.Println(err)
	}
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}

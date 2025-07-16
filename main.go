package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Contact struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ContactService struct {
	Contacts map[int]Contact
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (c *ContactService) Create(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := len(c.Contacts) + 1
	contact.Id = id

	c.Contacts[id] = contact

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contact)
	w.WriteHeader(http.StatusCreated)
}

func (c *ContactService) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contacts []Contact

	for _, ct := range c.Contacts {
		contacts = append(contacts, ct)
	}

	json.NewEncoder(w).Encode(contacts)
}

func (c *ContactService) Get(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	if val, ok := c.Contacts[id]; ok {
		json.NewEncoder(w).Encode(val)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func (c *ContactService) Delete(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	if _, ok := c.Contacts[id]; ok {
		delete(c.Contacts, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func (c *ContactService) Update(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	var contact Contact
	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := c.Contacts[id]; ok {
		contact.Id = id
		c.Contacts[id] = contact
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func HandleUpdateContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Update(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func handleGetContacts(w http.ResponseWriter, r *http.Request, service *ContactService) {
	w.Header().Set("Content-Type", "application/json")
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Get(w, r, id)
	} else {
		service.List(w, r)
	}
}

func handleCreateContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	service.Create(w, r)
}

func handleDeleteContact(w http.ResponseWriter, r *http.Request, service *ContactService) {
	q := r.URL.Query()
	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Delete(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func main() {
	service := &ContactService{Contacts: make(map[int]Contact)}
	mux := http.NewServeMux()

	mux.HandleFunc("/contacts", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateContact(w, r, service)
		case http.MethodGet:
			handleGetContacts(w, r, service)
		case http.MethodDelete:
			handleDeleteContact(w, r, service)
		case http.MethodPut:
			HandleUpdateContact(w, r, service)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}))
	fmt.Println("Server listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

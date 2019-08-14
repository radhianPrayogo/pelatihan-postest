package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"log"
	"net/http"
)

type Event struct {
	ID    int    `json:"ID"`
	Title string `json:"Title"`
	Date  string `json:"Date"`
	Place string `json:"Place"`
}

type Response struct {
	Error   bool   `json:"Error"`
	Message string `json:"Message"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/events", getEvents).Methods("GET")
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", selectEvent).Methods("GET")
	router.HandleFunc("/events/update/{id}", updateEvent).Methods("PATCH")
	log.Fatal(http.ListenAndServe(":9090", router))
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/postest_radhian")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Postest RESTful API")
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	var events []Event

	db := connect()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	results, err := db.Query("select * from events")
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	for results.Next() {
		var event Event

		err = results.Scan(&event.ID, &event.Title, &event.Date, &event.Place)
		if err != nil {
			panic(err.Error())
		}

		events = append(events, event)
	}

	json.NewEncoder(w).Encode(events)
}

func selectEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var event Event

	db := connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Query("select * from events where id = ?", eventID)
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	for result.Next() {
		err := result.Scan(&event.ID, &event.Title, &event.Date, &event.Place)
		if err != nil {
			log.Fatal(err)
		}
	}

	json.NewEncoder(w).Encode(event)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	db := connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	title := r.FormValue("title")
	date := r.FormValue("date")
	place := r.FormValue("place")

	update, err := db.Prepare("UPDATE events SET title=?, date=?, place=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	update.Exec(title, date, place, eventID)

	json.NewEncoder(w).Encode(Response{false, "Update data success"})
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	db := connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	title := r.FormValue("title")
	date := r.FormValue("date")
	place := r.FormValue("place")

	createEvent, err := db.Prepare("INSERT INTO events(title, date, place) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	createEvent.Exec(title, date, place)

	json.NewEncoder(w).Encode(Response{false, "Create data success"})
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	db := connect()
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	delEvent, err := db.Prepare("DELETE FROM events WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delEvent.Exec(eventID)
}

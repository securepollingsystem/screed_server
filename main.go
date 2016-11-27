package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/securepollingsystem/tallyspider/screed"
)

var Debug = false

func init() {
	if os.Getenv("DEBUG") == "1" {
		Debug = true
	}
}

func main() {
	db := mustInitDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/screed/{pubkey}", GetScreed(db)).Methods("GET")
	r.HandleFunc("/screed", PostScreed(db)).Methods("POST")

	http.Handle("/", r)

	listenAddr := ":8000"
	if port := os.Getenv("PORT"); port != "" {
		listenAddr = ":" + port
	}

	log.Printf("Listening on %v\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, r))
}

func GetScreed(db *bolt.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		pubkeyhex := mux.Vars(req)["pubkey"]
		screedB, err := GetScreedByPubkey(db, pubkeyhex)
		if err != nil {
			if err == ErrScreedNotFound {
				writeErrorStatus(w, err.Error(), http.StatusNotFound, err)
				return
			}
			writeError(w, "Error fetching screed ", err)
			return
		}
		w.Write(screedB)
		return
	}
}

func PostScreed(db *bolt.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			writeError(w, "Error reading POST data", err)
			return
		}
		defer req.Body.Close()

		screed, err := screed.DeserializeScreed(string(body))
		if err != nil {
			// TODO: Consider giving less info back to the user, for
			// security reasons
			writeErrorStatus(w, "Invalid POST: "+err.Error(),
				http.StatusBadRequest, nil)
			return
		}

		if err = screed.Valid(); err != nil {
			writeErrorStatus(w, "Invalid screed: "+err.Error(),
				http.StatusBadRequest, err)
		}

		pubkeyhex, err := CreateScreedByPubkey(db, screed)
		if err != nil {
			if err == ErrScreedExists {
				writeErrorStatus(w, err.Error(), http.StatusConflict, err)
				return
			}
			writeError(w, "Error saving new screed ", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(pubkeyhex))
	}
}

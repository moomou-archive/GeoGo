package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

const (
	POST    = "POST"
	GET     = "GET"
	DELETE  = "DELETE"
	PUT     = "PUT"
	OPTIONS = "OPTIONS"
)

type handlerWithError func(http.ResponseWriter, *http.Request) error
type handlerWithDB_and_Error func(http.ResponseWriter, *http.Request, *Trigger) error
type handler func(http.ResponseWriter, *http.Request)

type InvalidRequest struct {
	reason string
}

func (this *InvalidRequest) Error() string {
	return "Invalid request: " + this.reason
}

func getDBConnection() *sql.DB {
	dbName := os.Getenv("DB")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	pwd := os.Getenv("PASSWORD")
	user := os.Getenv("USER")

	if port == "" {
		port = "5432"
	}
	if host == "" {
		host = "localhost"
	}

	connectUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pwd, host, port, dbName)
	db, err := sql.Open("postgres", connectUrl)

	if err != nil {
		panic(err)
	}

	return db
}

func requestHandlerWithDB(db *sql.DB, methods map[string]handlerWithDB_and_Error) handler {
	t := newTrigger(db)

	return func(w http.ResponseWriter, r *http.Request) {
		// CORS header
		origin := r.Header["Origin"]

		if len(origin) == 1 {
			w.Header().Add("Access-Control-Allow-Origin", origin[0])
		}

		w.Header().Add("Access-Control-Allow-Headers",
			"DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,"+
				"If-Modified-Since,Cache-Control,Content-Type,Referer,x-access-token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		// By default, return json
		w.Header().Set("Content-Type", "application/json")

		if fn, ok := methods[r.Method]; ok {
			if err := fn(w, r, t); err != nil {
				switch err.(type) {
				case *InvalidRequest:
					{
						http.Error(w, err.Error(), http.StatusBadRequest)
					}
				default:
					{
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		} else if r.Method == OPTIONS {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Not supported", http.StatusNotFound)
		}
	}
}

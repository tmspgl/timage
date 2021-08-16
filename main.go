package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"timage.flomas.net/db"
	"timage.flomas.net/handler"
)

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	h := &handler.Handler{
		DB: dbConnection,
	}

	router:= mux.NewRouter()
    router.Use(commonMiddleware)
	router.Use(handlers.CORS())

	router.HandleFunc("/v1/image", h.RetrieveImages).Methods("GET")
    router.HandleFunc("/v1/image", h.CreateImage).Methods("PUT")

	router.HandleFunc("/v1/image/{imageId}", h.RetrieveImage).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
    allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
    allowedOrigins := handlers.AllowedOrigins([]string{"*"})

    srv := &http.Server{
    	Handler:      handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(loggedRouter),
    	Addr:         ":8080",
    	WriteTimeout: 15 * time.Second,
    	ReadTimeout:  15 * time.Second,
    }
	log.Fatal(srv.ListenAndServe())

}
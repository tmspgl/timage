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
	
	// Interval to check db for images to send.
	interval := 60 * time.Second
	h.StartImageFetch(interval)

	router := mux.NewRouter()
	router.Use(commonMiddleware)
	router.Use(handlers.CORS())

	// TODO currently serving last sent image because sending of images not implemented yet.
	router.HandleFunc("/", h.GetLatest).Methods("GET")

	router.HandleFunc("/v1/user", h.GetUser).Methods("GET")
	router.HandleFunc("/v1/user", h.CreateUser).Methods("POST")
	router.HandleFunc("/v1/user", h.EditUserInfo).Methods("PUT")

	router.HandleFunc("/v1/user/test", h.Test).Methods("GET")
	router.HandleFunc("/v1/user/test/images", h.Test2).Methods("GET")

	router.HandleFunc("/v1/user/{userId}/image", h.GetUserImage).Methods("GET")
	router.HandleFunc("/v1/user/{userId}/image", h.EditUserImage).Methods("PUT")

	router.HandleFunc("/v1/image", h.RetrieveAllImages).Methods("GET")
	router.HandleFunc("/v1/image", h.CreateImage).Methods("POST")

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

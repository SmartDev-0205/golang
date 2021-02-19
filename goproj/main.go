package main

import (
	"fmt"
	"goproj/controller"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("go server started in localhost:3000")
	r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/home", controller.HomeHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

}

package main

import (
	"anotherTZ/components/handlers"
	"anotherTZ/components/utils"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func main() {
	utils.MakeMigration()
	r := mux.NewRouter()
	r.HandleFunc("/info", handlers.GetPeople).Methods("GET")
	r.HandleFunc("/info", handlers.SetPeople).Methods("POST")
	r.HandleFunc("/info", handlers.UpdatePeople).Methods("PUT")
	r.HandleFunc("/info", handlers.DeletePeople).Methods("DELETE")

	r.HandleFunc("/task", handlers.GetTask).Methods("GET")
	r.HandleFunc("/task", handlers.SetTask).Methods("POST")
	r.HandleFunc("/task", handlers.CloseTask).Methods("PUT")

	r.HandleFunc("/swagger", httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", r)
}

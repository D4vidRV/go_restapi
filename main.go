package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/D4vidRV/go_restapi/db"
	"github.com/D4vidRV/go_restapi/models"
	"github.com/D4vidRV/go_restapi/routes"
)

func main() {
	db.NewDB()

	db.DB().AutoMigrate(
		&models.User{},
		&models.Task{},
	)

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Users routes
	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods("DELETE")

	// Tasks routes
	r.HandleFunc("/tasks", routes.GetTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", routes.GetTaskHandler).Methods("GET")
	r.HandleFunc("/tasks", routes.PostTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", routes.DeleteTaskHandler).Methods("DELETE")

	http.ListenAndServe(":3000", r)
}

package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/D4vidRV/go_restapi/db"
	"github.com/D4vidRV/go_restapi/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	var resp *models.HttpResponse
	w.Header().Set("Content-Type", "application/json")

	if err := db.DB().Find(&tasks).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = &models.HttpResponse{Message: "Internal Server Error"}

		resJSON, _ := json.Marshal(resp)
		w.Write(resJSON)
		return
	}

	if len(tasks) == 0 {
		resp = &models.HttpResponse{Message: "OK", Data: []models.Task{}}
		log.Println(resp)
	} else {
		tasksJSON, _ := json.Marshal(tasks)
		resp = &models.HttpResponse{Message: "OK", Data: tasksJSON}
	}

	resJSON, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var resp *models.HttpResponse
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	if err := db.DB().First(&task, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			resp = &models.HttpResponse{Message: "Task not found"}
			respJSON, _ := json.Marshal(resp)
			w.Write(respJSON)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		resp = &models.HttpResponse{Message: fmt.Sprintf("Internal Server Error: %v", err.Error())}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	respJSON, _ := json.Marshal(task)
	w.Write(respJSON)
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var resp *models.HttpResponse
	w.Header().Set("Content-Type", "application/json")

	// Parsing of the request
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatalf("Error al hacer el decode json %v", err)
	}

	// Make validation of dto
	validate := validator.New()
	err = validate.Struct(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp = &models.HttpResponse{Message: fmt.Sprintf("Bad Request: %v", err.Error())}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	// Create a Post
	if err := db.DB().Create(&task).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = &models.HttpResponse{Message: fmt.Sprintf("Internal Server Error %v", err.Error())}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	// Return respornse in JSON format
	respJSON, _ := json.Marshal(task)
	w.Write(respJSON)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	params := mux.Vars(r)
	var resp *models.HttpResponse
	w.Header().Set("Content-Type", "application/json")

	// Find task
	db.DB().First(&task, params["id"])
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp = &models.HttpResponse{Message: fmt.Sprintf("Task with id %v not exist", params["id"])}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	// If found task, delete it
	db.DB().Delete(&task)
	w.WriteHeader(http.StatusOK)
	resp = &models.HttpResponse{Message: fmt.Sprintf("Task with id %v deleted", params["id"])}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

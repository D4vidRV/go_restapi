package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/D4vidRV/go_restapi/db"
	"github.com/D4vidRV/go_restapi/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	db.DB().First(&task, params["id"])
	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	json.NewEncoder(w).Encode(&task)
}

func PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	w.Header().Set("Content-Type", "application/json")
	// Parsin de la request
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Fatalf("Error al hacer el decode json %v", err)
	}

	// Hacer  validacion del dto
	validate := validator.New()
	err = validate.Struct(task)
	if err != nil {
		errMsg := fmt.Sprintf("Error de dto: %v", err)
		resp := map[string]interface{}{"message": errMsg}
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

	createdTask := db.DB().Create(&task)
	if createdTask.Error != nil {
		errMsg := fmt.Sprintf("Error al crear tarea: %v", createdTask.Error)
		resp := map[string]interface{}{"message": errMsg}
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al convertir a JSON"))
		return
	}

	w.Write(taskJSON)
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
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return
	}

	// If found task, delete it
	db.DB().Delete(&task)
	w.WriteHeader(http.StatusOK)
	resp = &models.HttpResponse{Message: fmt.Sprintf("Task with id %v deleted", params["id"])}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

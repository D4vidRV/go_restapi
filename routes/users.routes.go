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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	var resp *models.HttpResponse
	w.Header().Set("Content-Type", "application/json")

	if err := db.DB().Find(&users).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = &models.HttpResponse{Message: "Internal Server Error"}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	respJSON, _ := json.Marshal(users)

	w.Write(respJSON)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var resp *models.HttpResponse
	var user models.User
	w.Header().Set("Content-Type", "application/json")

	if err := db.DB().First(&user, params["id"]).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			resp = &models.HttpResponse{Message: "User not found"}
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

	// Add Tasks to user using FK
	db.DB().Model(&user).Association("Tasks").Find(&user.Tasks)
	respJSON, _ := json.Marshal(user)
	w.Write(respJSON)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var resp *models.HttpResponse
	w.Header().Set("Content-Type", "application/json")

	// Parsin de la request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Error al hacer el decode json %v", err)
	}

	// Make validation of DTO
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		resp = &models.HttpResponse{Message: fmt.Sprintf("Internal Server Error: %v", err.Error())}
		respJSON, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respJSON)
		return
	}

	if err := db.DB().Create(&user).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp = &models.HttpResponse{Message: fmt.Sprintf("Internal Server Error %v", err.Error())}
		respJSON, _ := json.Marshal(resp)
		w.Write(respJSON)
		return
	}

	// Return the response in JSON format
	respJSON, _ := json.Marshal(user)
	w.Write(respJSON)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var resp *models.HttpResponse

	params := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	db.DB().First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		resp = &models.HttpResponse{Message: fmt.Sprintf("User with id %v not exist", params["id"])}
		jsonResp, _ := json.Marshal(resp)
		w.Write(jsonResp)
		return
	}

	db.DB().Delete(&user)
	w.WriteHeader(http.StatusOK)
	resp = &models.HttpResponse{Message: fmt.Sprintf("User with id %v deleted", params["id"])}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

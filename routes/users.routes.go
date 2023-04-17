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

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	db.DB().Find(&users)
	w.Header().Set("Content-Type", "application/json")

	usersJSON, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al convertir a JSON"))
		return
	}

	w.Write(usersJSON)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User
	w.Header().Set("Content-Type", "application/json")

	db.DB().First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))

		return
	}

	db.DB().Model(&user).Association("Tasks").Find(&user.Tasks)

	json.NewEncoder(w).Encode(&user)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	w.Header().Set("Content-Type", "application/json")
	// Parsin de la request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Error al hacer el decode json %v", err)
	}

	// Hacer  validacion del dto
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		errMsg := fmt.Sprintf("Error de dto: %v", err)
		resp := map[string]interface{}{"message": errMsg}
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

	createdUser := db.DB().Create(&user)
	if createdUser.Error != nil {
		errMsg := fmt.Sprintf("Error al crear usuario: %v", createdUser.Error)
		resp := map[string]interface{}{"message": errMsg}
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
		return
	}

	// Devolver la respuesta en formato JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error al convertir a JSON"))
		return
	}

	w.Write(userJSON)

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

package handlers

import (
	"Golang-DeFiExplorer/internal/db" 
	"Golang-DeFiExplorer/internal/models"
	"Golang-DeFiExplorer/internal/repository"
	"encoding/json"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при декодировании запроса: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Все поля обязательны", http.StatusBadRequest)
		return
	}

	dbInstance := db.GetDBInstance() 

	if err := repository.SaveUser(dbInstance, &user); err != nil { 
		http.Error(w, "Ошибка при сохранении пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

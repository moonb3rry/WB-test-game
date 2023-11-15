package http

import (
	"WB_game/internal/middleware"
	"WB_game/internal/model"
	"encoding/json"
	"net/http"
)

func (c *controller) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req model.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	err = c.userService.RegisterUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func (c *controller) LogInUserHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	token, err := c.userService.LogInUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c *controller) AboutHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaims(r)
	if err != nil {
		http.Error(w, "Ошибка извлечения claims: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.userService.AboutUser(r.Context(), claims)
	if err != nil {
		http.Error(w, "Ошибка получения данных пользователя: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

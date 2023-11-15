package http

import (
	"WB_game/internal/middleware"
	"encoding/json"
	"net/http"
)

func (c *controller) TasksHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaims(r) // Функция для извлечения claims из JWT
	if err != nil {
		http.Error(w, "Ошибка извлечения claims: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := c.taskService.GetAvailableTasks(r.Context(), *claims)
	if err != nil {
		http.Error(w, "Ошибка получения данных пользователя: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func (c *controller) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := c.taskService.GenerateTasks(r.Context())
	if err != nil {
		http.Error(w, "Ошибка при генерации заданий: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&tasks)
}

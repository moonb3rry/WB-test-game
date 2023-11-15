package http

import (
	"WB_game/internal/middleware"
	"WB_game/internal/model"
	"encoding/json"
	"net/http"
)

func (c *controller) StartGameHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := middleware.ExtractClaims(r) // Функция для извлечения claims из JWT
	if err != nil {
		http.Error(w, "Ошибка извлечения claims: "+err.Error(), http.StatusBadRequest)
		return
	}
	var start model.StartTask
	err = json.NewDecoder(r.Body).Decode(&start)
	if err != nil {
		return
	}
	result, err := c.gameService.StartGame(r.Context(), claims, start)
	if err != nil {
		http.Error(w, "Ошибка получения данных пользователя: "+err.Error(), http.StatusBadRequest)
		return
	}
	if result == true {
		json.NewEncoder(w).Encode(map[string]string{"message": "Ура! Победа"})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Вы проиграли"})
	}

}

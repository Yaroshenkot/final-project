package api

import (
	"encoding/json"
	"net/http"
	"os"
)

// signinHandler обрабатывает POST-запрос /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "метод не поддерживается"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if req.Password != expectedPassword {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	token, err := GenerateToken()
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка генерации токена"})
		return
	}

	writeJSON(w, map[string]string{"token": token})
}

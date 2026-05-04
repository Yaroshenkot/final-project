package api

import (
	"encoding/json"
	"net/http"

	"final-project/pkg/db"
)

// getTaskHandler обрабатывает GET-запрос /api/task?id=N
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if task == nil {
		writeJSON(w, map[string]string{"error": "задача не найдена"})
		return
	}

	writeJSON(w, task)
}

// updateTaskHandler обрабатывает PUT-запрос /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	// Проверяем ID
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	// Проверяем заголовок
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверяем дату и правило повторения
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновляем задачу
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

package api

import (
	"net/http"
)

// Init регистрирует все API-обработчики
func Init() {
	// Публичные эндпоинты
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/signin", signinHandler)

	// Защищённые эндпоинты
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}

// taskHandler маршрутизирует запросы по HTTP-методу
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, map[string]string{"error": "метод не поддерживается"})
	}
}

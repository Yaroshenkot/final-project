package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"final-project/pkg/db"
)

// addTaskHandler обрабатывает POST-запрос на добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var task db.Task

	// 1. Десериализуем JSON
	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	// 2. Проверяем заголовок
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// 3. Проверяем дату и правило повторения
	if err = checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 4. Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// 5. Возвращаем ID
	writeJSON(w, http.StatusOK, map[string]string{"id": toString(id)})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана — берём сегодня
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	// Парсим дату
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return errors.New("дата представлена в формате, отличном от 20060102")
	}

	// Если задача с повторением — вычисляем следующую дату
	var nextDate string
	if task.Repeat != "" {
		nextDate, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если сегодня больше даты задачи
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = nextDate
		}
	}

	return nil
}

// writeJSON записывает данные в формате JSON
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "ошибка формирования ответа", http.StatusInternalServerError)
	}
}

// toString преобразует int64 в строку
func toString(id int64) string {
	return fmt.Sprintf("%d", id)
}

package api

import (
	"net/http"

	"final-project/pkg/db"
)

const defaultLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.SearchTasks(search, defaultLimit)
	} else {
		tasks, err = db.Tasks(defaultLimit)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
}

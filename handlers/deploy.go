package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"
	"strings"
	"deploy-server/config"
	"deploy-server/models"
)

func DeployHandler(db *sql.DB, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := r.URL.Query().Get("id")
		if appID == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		var app models.Application
		err := db.QueryRow("SELECT id, name, commands, path FROM applications WHERE id = ?", appID).Scan(&app.ID, &app.Name, &app.Commands, &app.Path)
		if err != nil {
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}

		session := models.DeploymentSession{
			AppID:     app.ID,
			StartTime: time.Now(),
			Status:    "started",
		}

		res, err := db.Exec("INSERT INTO deployment_session (app_id, start_time, status) VALUES (?, ?, ?)", session.AppID, session.StartTime, session.Status)
		if err != nil {
			http.Error(w, "Failed to create deployment session", http.StatusInternalServerError)
			return
		}

		sessionID, _ := res.LastInsertId()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int64{"session_id": sessionID})

		go executeCommands(db, sessionID, app.Commands, cfg, app.Path)
	}
}

func executeCommands(db *sql.DB, sessionID int64, commands string, cfg config.Config, appPath string) {
	cmdList := strings.Split(commands, "\n")
	for _, cmd := range cmdList {
		step := models.DeploymentStep{
			SessionID:  sessionID,
			Command:    cmd,
			ExecutedAt: time.Now(),
		}

		// Создаем команду с указанием рабочей директории
		command := exec.Command(cfg.Shell.Binary, append(cfg.Shell.Args, cmd)...)
		command.Dir = appPath // Устанавливаем рабочую директорию

		out, err := command.CombinedOutput()

		if err != nil {
			step.Status = "failed"
			step.Output = err.Error()
		} else {
			step.Status = "success"
			step.Output = string(out)
		}

		db.Exec("INSERT INTO deployment_steps (session_id, command, output, status, executed_at) VALUES (?, ?, ?, ?, ?)",
			step.SessionID, step.Command, step.Output, step.Status, step.ExecutedAt)
	}

	db.Exec("UPDATE deployment_session SET end_time = ?, status = 'completed' WHERE id = ?", time.Now(), sessionID)
}

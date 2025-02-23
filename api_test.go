package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"deploy-server/config"
	"deploy-server/models"
	"deploy-server/utils"
	"deploy-server/handlers"

	_ "github.com/mattn/go-sqlite3" // Импорт драйвера SQLite
)

// TestMain настраивает тестовое окружение.
func TestMain(m *testing.M) {
	// Инициализация тестовой базы данных
	dbPath := "file::memory:?cache=shared"
	db, err := models.InitDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Добавляем тестовое приложение
	_, err = db.Exec("INSERT INTO applications (name, commands, path) VALUES (?, ?, ?)", "Test App", "echo 'Hello, World!'", "/tmp/test-app")
	if err != nil {
		panic(err)
	}

	// Запуск тестов
	code := m.Run()

	// Очистка после тестов
	os.Remove(dbPath)
	os.Exit(code)
}

// TestDeployHandler проверяет обработчик /deploy.
func TestDeployHandler(t *testing.T) {
	// Загружаем конфигурацию
	cfg := config.Config{
		Shell: struct {
			Binary string   `yaml:"binary"`
			Args   []string `yaml:"args"`
		}{
			Binary: "sh",
			Args:   []string{"-c"},
		},
	}

	// Инициализация базы данных
	dbPath := "file::memory:?cache=shared"
	db, err := models.InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер
	handler := http.HandlerFunc(utils.BasicAuth(handlers.DeployHandler(db, cfg)))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", ts.URL+"/deploy?id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Устанавливаем Basic Auth
	req.SetBasicAuth("admin", "password")

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Проверяем тело ответа
	var result map[string]int64
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	sessionID, ok := result["session_id"]
	if !ok {
		t.Error("Response does not contain session_id")
	}

	// Проверяем, что сессия деплоя создана в базе данных
	var session models.DeploymentSession
	err = db.QueryRow("SELECT id, app_id, status FROM deployment_session WHERE id = ?", sessionID).Scan(&session.ID, &session.AppID, &session.Status)
	if err != nil {
		t.Fatalf("Failed to retrieve deployment session: %v", err)
	}

	if session.AppID != 1 || session.Status != "started" {
		t.Errorf("Deployment session does not match. Got: %+v", session)
	}
}

// TestDeployHandlerUnauthorized проверяет обработчик /deploy без авторизации.
func TestDeployHandlerUnauthorized(t *testing.T) {
	// Загружаем конфигурацию
	cfg := config.Config{
		Shell: struct {
			Binary string   `yaml:"binary"`
			Args   []string `yaml:"args"`
		}{
			Binary: "sh",
			Args:   []string{"-c"},
		},
	}

	// Инициализация базы данных
	dbPath := "file::memory:?cache=shared"
	db, err := models.InitDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создаем тестовый сервер
	handler := http.HandlerFunc(utils.BasicAuth(handlers.DeployHandler(db, cfg)))
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Создаем HTTP-запрос без авторизации
	req, err := http.NewRequest("POST", ts.URL+"/deploy?id=1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code 401, got %d", resp.StatusCode)
	}
}
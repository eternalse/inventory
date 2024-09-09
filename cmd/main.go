package main

import (
	"database/sql"
	"fmt"
	"inva/config"
	"inva/pkg/logging"
	"inva/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Получение текущего рабочего каталога
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении рабочего каталога:", err)
		return
	}

	fmt.Println("Текущий рабочий каталог:", cwd)

	// Указываем путь к конфигурационному файлу
	configPath := "/home/alena/inva/config/config.yaml"

	fmt.Printf("Используется конфигурационный файл: %s\n", configPath)

	// Загружаем конфигурацию
	appConfig, err := config.LoadAppConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	// Настраиваем логирование
	logger, err := logging.NewLogger(appConfig.Logging.Level)
	if err != nil {
		log.Fatalf("Ошибка при создании логгера: %v", err)
	}

	// Подключение к базе данных
	db, err := sql.Open("postgres", config.GetDatabaseURL(appConfig))
	if err != nil {
		logger.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		logger.Fatalf("Ошибка при проверке соединения с базой данных: %v", err)
	}

	// Настройка маршрутизации
	r := mux.NewRouter()
	routes.SetupRoutes(r, db)

	// Запуск сервера
	port := appConfig.Server.Port
	if port[0] == ':' {
		port = port[1:] // Убираем двоеточие, если оно есть
	}
	serverAddress := fmt.Sprintf(":%s", port)
	fmt.Println("Сервер работает на порту", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
	logger.Info("Приложение запущено")
}

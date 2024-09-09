package routes

import (
	"database/sql"
	"inva/handlers"
	"inva/services"

	"github.com/gorilla/mux"
)

// SetupRoutes конфигурирует маршруты и обработчики
func SetupRoutes(r *mux.Router, db *sql.DB) {
	// Создание сервисов
	employeeService := services.NewEmployeeService(db)
	equipmentService := services.NewEquipmentService(db)

	// Создание обработчиков с передачей сервисов
	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	equipmentHandler := handlers.NewEquipmentHandler(equipmentService)

	// Маршруты для сотрудников
	r.HandleFunc("/employees", employeeHandler.GetAllEmployeesHandler).Methods("GET")
	r.HandleFunc("/employees", employeeHandler.CreateEmployeeHandler).Methods("POST")
	r.HandleFunc("/employees/{id:[0-9]+}", employeeHandler.GetEmployeeHandler).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", employeeHandler.UpdateEmployeeHandler).Methods("PUT")
	r.HandleFunc("/employees/{id:[0-9]+}", employeeHandler.DeleteEmployeeHandler).Methods("DELETE")

	// Маршруты для оборудования
	r.HandleFunc("/equipment", equipmentHandler.GetAllEquipmentHandler).Methods("GET")
	r.HandleFunc("/equipment", equipmentHandler.CreateEquipmentHandler).Methods("POST")
	r.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.GetEquipmentHandler).Methods("GET")
	r.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.UpdateEquipmentHandler).Methods("PUT")
	r.HandleFunc("/equipment/{id:[0-9]+}", equipmentHandler.DeleteEquipmentHandler).Methods("DELETE")

	// Назначение оборудования пользователю
	r.HandleFunc("/equipment/{equipment_id:[0-9]+}/assign/user/{user_id:[0-9]+}", equipmentHandler.AssignEquipmentToUser).Methods("POST")

	// Возврат оборудования
	r.HandleFunc("/equipment/{id:[0-9]+}/return", equipmentHandler.ReturnEquipmentHandler).Methods("PUT")

	// Детали оборудования
	r.HandleFunc("/equipment/{id:[0-9]+}/details", equipmentHandler.GetEquipmentDetailsHandler).Methods("GET")
}

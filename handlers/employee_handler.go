package handlers

import (
	"encoding/json"
	"inva/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// EmployeeHandler представляет обработчик для операций с сотрудниками
type EmployeeHandler struct {
	service *services.EmployeeService
}

// NewEmployeeHandler создаёт новый экземпляр EmployeeHandler
func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

// CreateEmployeeHandler обрабатывает HTTP запрос для создания нового сотрудника
func (h *EmployeeHandler) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var employee services.Employee

	// Декодируем тело запроса в структуру Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logrus.WithError(err).Error("Ошибка при декодировании запроса")
		return
	}

	// Создаем сотрудника через сервис
	createdEmployee, err := h.service.CreateEmployee(&employee)
	if err != nil {
		http.Error(w, "Error creating employee", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при создании сотрудника")
		return
	}

	// Возвращаем ответ с кодом 201 (Created) и данными о созданном сотруднике
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdEmployee); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при кодировании ответа")
	}
}

// GetAllEmployeesHandler обрабатывает HTTP запрос для получения всех сотрудников
func (h *EmployeeHandler) GetAllEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем всех сотрудников через сервис
	employees, err := h.service.GetAllEmployees()
	if err != nil {
		http.Error(w, "Error retrieving employees", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при получении списка сотрудников")
		return
	}

	// Возвращаем ответ с кодом 200 (OK) и данными о сотрудниках
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при кодировании ответа")
	}
}

// UpdateEmployeeHandler обрабатывает HTTP запрос для обновления данных сотрудника
func (h *EmployeeHandler) UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		logrus.WithError(err).Error("Ошибка при преобразовании ID")
		return
	}

	var employee services.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logrus.WithError(err).Error("Ошибка при декодировании запроса")
		return
	}

	// Обновляем сотрудника через сервис
	err = h.service.UpdateEmployee(id, employee.Name)
	if err != nil {
		http.Error(w, "Error updating employee", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при обновлении сотрудника")
		return
	}

	// Возвращаем ответ с кодом 204 (No Content) после успешного обновления
	w.WriteHeader(http.StatusNoContent)
}

// DeleteEmployeeHandler обрабатывает HTTP запрос для удаления сотрудника
func (h *EmployeeHandler) DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		logrus.WithError(err).Error("Ошибка при преобразовании ID сотрудника")
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		http.Error(w, "Error deleting employee", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при удалении сотрудника")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetEmployeeHandler обрабатывает HTTP запрос для получения одного сотрудника по ID
func (h *EmployeeHandler) GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		logrus.WithError(err).Error("Ошибка при преобразовании ID")
		return
	}

	employee, err := h.service.GetEmployeeByID(id)
	if err != nil {
		http.Error(w, "Error retrieving employee", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при получении сотрудника")
		return
	}

	if employee == nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Возвращаем ответ с кодом 200 (OK) и данными о сотруднике
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		logrus.WithError(err).Error("Ошибка при кодировании ответа")
	}
}

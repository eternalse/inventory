package utils

import (
	"encoding/json"
	"inva/models"
	"inva/services"
	"net/http"
)

// RespondWithJSON отправляет JSON-ответ с указанным статусом и данными
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RespondWithError отправляет ошибку в формате JSON
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": message})
}

// ParseRequestBody парсит тело запроса в указанную структуру
func ParseRequestBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GenerateID создаёт уникальный идентификатор
func GenerateID() string {
	return "unique-id-placeholder"
}

// Преобразование из *models.Employee в *services.Employee
func ConvertModelToServiceEmployee(modelEmployee *models.Employee) *services.Employee {
	return &services.Employee{
		ID:   modelEmployee.ID,
		Name: modelEmployee.Name,
	}
}

// Преобразование из *services.Employee в *models.Employee
func ConvertServiceToModelEmployee(serviceEmployee *services.Employee) *models.Employee {
	return &models.Employee{
		ID:   serviceEmployee.ID,
		Name: serviceEmployee.Name,
	}
}

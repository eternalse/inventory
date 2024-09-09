package handlers

import (
	"encoding/json"
	"inva/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// EquipmentHandler представляет обработчик для работы с оборудованием
type EquipmentHandler struct {
	service *services.EquipmentService
}

// NewEquipmentHandler создаёт новый экземпляр EquipmentHandler
func NewEquipmentHandler(service *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

// AssignEquipmentToUser закрепляет оборудование за пользователем
func (h *EquipmentHandler) AssignEquipmentToUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	equipmentIDStr := vars["equipment_id"]
	userIDStr := vars["user_id"]

	// Преобразуем ID оборудования и пользователя в целое число
	equipmentID, err := strconv.Atoi(equipmentIDStr)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentIDStr,
			"user_id":      userIDStr,
		}).Error("Ошибка при преобразовании ID оборудования")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"user_id":      userIDStr,
			"equipment_id": equipmentIDStr,
		}).Error("Ошибка при преобразовании ID пользователя")
		return
	}

	// Присваиваем оборудование пользователю
	err = h.service.AssignEquipmentToUser(equipmentID, userID)
	if err != nil {
		http.Error(w, "Failed to assign equipment to user: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentID,
			"user_id":      userID,
		}).Error("Ошибка при назначении оборудования пользователю")
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusNoContent)
	logrus.WithFields(logrus.Fields{
		"equipment_id": equipmentID,
		"user_id":      userID,
	}).Info("Оборудование успешно назначено пользователю")
}

// ReturnEquipmentHandler обрабатывает запрос на возврат оборудования
func (h *EquipmentHandler) ReturnEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	equipmentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": vars["id"],
		}).Error("Ошибка при преобразовании ID оборудования")
		return
	}

	// Возвращаем оборудование от пользователя
	if err := h.service.ReturnEquipmentFromUser(equipmentID); err != nil {
		http.Error(w, "Error returning equipment: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentID,
		}).Error("Ошибка при возврате оборудования")
		return
	}

	w.WriteHeader(http.StatusOK) // Статус успешного выполнения
	logrus.WithField("equipment_id", equipmentID).Info("Оборудование успешно возвращено")
}

// GetEquipmentDetailsHandler возвращает информацию о конкретном оборудовании
func (h *EquipmentHandler) GetEquipmentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	equipmentIDStr := mux.Vars(r)["id"]

	equipmentID, err := strconv.Atoi(equipmentIDStr)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentIDStr,
		}).Error("Ошибка при преобразовании ID оборудования")
		return
	}

	// Получаем детали оборудования
	equipment, err := h.service.GetEquipmentDetails(equipmentID)
	if err != nil {
		http.Error(w, "Error retrieving equipment details: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentID,
		}).Error("Ошибка при получении деталей оборудования")
		return
	}

	if equipment == nil {
		http.NotFound(w, r)
		logrus.WithField("equipment_id", equipmentID).Warn("Оборудование не найдено")
		return
	}

	// Отправляем данные в JSON формате
	response, err := json.Marshal(equipment)
	if err != nil {
		http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": equipmentID,
		}).Error("Ошибка при кодировании ответа")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	logrus.WithField("equipment_id", equipmentID).Info("Детали оборудования успешно возвращены")
}

// CreateEquipmentHandler обрабатывает создание нового оборудования
func (h *EquipmentHandler) CreateEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	var equipment services.Equipment
	// Декодируем JSON данные
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  r.Body,
		}).Error("Ошибка при декодировании запроса на создание оборудования")
		return
	}

	// Создаем новое оборудование
	createdEquipment, err := h.service.CreateEquipment(equipment.Model, equipment.SerialNumber, equipment.Status)
	if err != nil {
		http.Error(w, "Error creating equipment: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"equipment": equipment,
		}).Error("Ошибка при создании оборудования")
		return
	}

	// Отправляем успешный ответ с кодом 201 Created
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdEquipment); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":             err,
			"created_equipment": createdEquipment,
		}).Error("Ошибка при кодировании ответа")
	}
	logrus.WithField("equipment_id", createdEquipment.ID).Info("Оборудование успешно создано")
}

// GetEquipmentHandler обрабатывает получение оборудования по ID
func (h *EquipmentHandler) GetEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithField("error", "ID not found in URL").Error("Ошибка при получении ID оборудования из URL")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": idStr,
		}).Error("Ошибка при преобразовании ID оборудования")
		return
	}

	// Получаем оборудование по ID
	equipment, err := h.service.GetEquipmentByID(id)
	if err != nil {
		http.Error(w, "Error retrieving equipment: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": id,
		}).Error("Ошибка при получении оборудования")
		return
	}
	if equipment == nil {
		http.Error(w, "Equipment not found", http.StatusNotFound)
		logrus.WithField("equipment_id", id).Warn("Оборудование не найдено")
		return
	}

	// Отправляем успешный ответ с кодом 200 OK
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(equipment); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"equipment": equipment,
		}).Error("Ошибка при кодировании ответа")
	}
	logrus.WithField("equipment_id", id).Info("Оборудование успешно возвращено")
}

// GetAllEquipmentHandler обрабатывает получение списка всего оборудования
func (h *EquipmentHandler) GetAllEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем список оборудования
	equipmentList, err := h.service.GetAllEquipment()
	if err != nil {
		http.Error(w, "Error retrieving equipment list: "+err.Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("Ошибка при получении списка оборудования")
		return
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(equipmentList); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		logrus.WithField("error", err).Error("Ошибка при кодировании ответа")
	}
	logrus.Info("Список всего оборудования успешно возвращен")
}

// UpdateEquipmentHandler обрабатывает обновление данных оборудования
func (h *EquipmentHandler) UpdateEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	var equipment services.Equipment
	// Декодируем JSON данные
	if err := json.NewDecoder(r.Body).Decode(&equipment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error": err,
			"body":  r.Body,
		}).Error("Ошибка при декодировании запроса на обновление оборудования")
		return
	}

	// Обновляем оборудование
	err := h.service.UpdateEquipment(equipment.ID, equipment.Model, equipment.SerialNumber, equipment.Status)
	if err != nil {
		http.Error(w, "Error updating equipment: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":     err,
			"equipment": equipment,
		}).Error("Ошибка при обновлении оборудования")
		return
	}

	// Возвращаем успешный статус без контента
	w.WriteHeader(http.StatusNoContent)
	logrus.WithField("equipment_id", equipment.ID).Info("Оборудование успешно обновлено")
}

// DeleteEquipmentHandler обрабатывает удаление оборудования по ID
func (h *EquipmentHandler) DeleteEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithField("error", "ID not found in URL").Error("Ошибка при получении ID оборудования из URL")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid equipment ID", http.StatusBadRequest)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": idStr,
		}).Error("Ошибка при преобразовании ID оборудования")
		return
	}

	// Удаляем оборудование
	err = h.service.DeleteEquipment(id)
	if err != nil {
		http.Error(w, "Error deleting equipment: "+err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"equipment_id": id,
		}).Error("Ошибка при удалении оборудования")
		return
	}

	// Возвращаем успешный статус
	w.WriteHeader(http.StatusNoContent)
	logrus.WithField("equipment_id", id).Info("Оборудование успешно удалено")
}

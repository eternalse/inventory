package repositories

import (
	"database/sql"
	"fmt"
	"inva/models"
)

// SQLRepository реализует интерфейс EquipmentRepository
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository создает новый репозиторий для работы с оборудованием
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// CreateEquipment создает новую единицу оборудования
func (r *SQLRepository) CreateEquipment(equipment *models.Equipment) (*models.Equipment, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO equipment (model, status) VALUES ($1, $2) RETURNING id", equipment.Model, equipment.Status).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании оборудования: %v", err)
	}
	equipment.ID = id
	return equipment, nil
}

// GetEquipmentByID возвращает оборудование по его идентификатору
func (r *SQLRepository) GetEquipmentByID(id int) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.QueryRow("SELECT id, model, status, assigned_to FROM equipment WHERE id = $1", id).
		Scan(&equipment.ID, &equipment.Model, &equipment.Status, &equipment.AssignedTo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("оборудование с id %d не найдено", id)
		}
		return nil, fmt.Errorf("ошибка при получении оборудования: %v", err)
	}
	return &equipment, nil
}

// GetAllEquipment возвращает список всех единиц оборудования
func (r *SQLRepository) GetAllEquipment() ([]models.Equipment, error) {
	rows, err := r.db.Query("SELECT id, model, status, assigned_to FROM equipment")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех единиц оборудования: %v", err)
	}
	defer rows.Close()

	var equipmentList []models.Equipment
	for rows.Next() {
		var equipment models.Equipment
		if err := rows.Scan(&equipment.ID, &equipment.Model, &equipment.Status, &equipment.AssignedTo); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}
		equipmentList = append(equipmentList, equipment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе строк: %v", err)
	}
	return equipmentList, nil
}

// UpdateEquipment обновляет данные оборудования
func (r *SQLRepository) UpdateEquipment(equipment *models.Equipment) error {
	_, err := r.db.Exec("UPDATE equipment SET model = $1, status = $2 WHERE id = $3", equipment.Model, equipment.Status, equipment.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении оборудования: %v", err)
	}
	return nil
}

// DeleteEquipment удаляет оборудование из базы данных
func (r *SQLRepository) DeleteEquipment(id int) error {
	_, err := r.db.Exec("DELETE FROM equipment WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении оборудования: %v", err)
	}
	return nil
}

// AssignEquipmentToUser закрепляет оборудование за пользователем
func (r *SQLRepository) AssignEquipmentToUser(equipmentID, userID int) error {
	_, err := r.db.Exec("UPDATE equipment SET assigned_to = $1 WHERE id = $2", userID, equipmentID)
	if err != nil {
		return fmt.Errorf("ошибка при закреплении оборудования: %v", err)
	}
	return nil
}

// ReturnEquipmentFromUser возвращает оборудование обратно
func (r *SQLRepository) ReturnEquipmentFromUser(equipmentID int) error {
	_, err := r.db.Exec("UPDATE equipment SET assigned_to = NULL WHERE id = $1", equipmentID)
	if err != nil {
		return fmt.Errorf("ошибка при возврате оборудования: %v", err)
	}
	return nil
}

// GetEquipmentDetails возвращает подробности об оборудовании, включая информацию о том, за кем оно закреплено
func (r *SQLRepository) GetEquipmentDetails(id int) (*models.Equipment, error) {
	var equipment models.Equipment
	err := r.db.QueryRow("SELECT id, model, status, assigned_to FROM equipment WHERE id = $1", id).
		Scan(&equipment.ID, &equipment.Model, &equipment.Status, &equipment.AssignedTo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("оборудование с id %d не найдено", id)
		}
		return nil, fmt.Errorf("ошибка при получении оборудования: %v", err)
	}
	return &equipment, nil
}

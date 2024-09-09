package services

import (
	"database/sql"
	"fmt"
	"log"
)

// Equipment представляет модель оборудования
type Equipment struct {
	ID           int    `json:"id"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	AssignedTo   *int   `json:"assigned_to"`
	UpdatedAt    string `json:"updated_at"`
}

// EquipmentService представляет сервис для работы с оборудованием
type EquipmentService struct {
	db *sql.DB
}

// NewEquipmentService создаёт новый экземпляр EquipmentService
func NewEquipmentService(db *sql.DB) *EquipmentService {
	return &EquipmentService{db: db}
}

// AssignEquipmentToUser закрепляет оборудование за пользователем
func (s *EquipmentService) AssignEquipmentToUser(equipmentID, userID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %v", err)
	}

	_, err = tx.Exec("UPDATE equipment SET assigned_to = $1 WHERE id = $2", userID, equipmentID)
	if err != nil {
		tx.Rollback() // откат в случае ошибки
		return fmt.Errorf("ошибка при закреплении оборудования: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ошибка при подтверждении транзакции: %v", err)
	}

	return nil
}

// ReturnEquipmentFromUser возвращает оборудование обратно
func (s *EquipmentService) ReturnEquipmentFromUser(equipmentID int) error {
	_, err := s.db.Exec("UPDATE equipment SET assigned_to = NULL WHERE id = $1", equipmentID)
	if err != nil {
		return fmt.Errorf("ошибка при возврате оборудования: %v", err)
	}

	return nil
}

// GetEquipmentDetails возвращает подробности об оборудовании, включая информацию о том, за кем оно закреплено
func (s *EquipmentService) GetEquipmentDetails(id int) (*Equipment, error) {
	var equipment Equipment
	err := s.db.QueryRow(
		"SELECT id, model, status, assigned_to FROM equipment WHERE id = $1", id,
	).Scan(&equipment.ID, &equipment.Model, &equipment.Status, &equipment.AssignedTo)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("оборудование с id %d не найдено", id)
		}
		return nil, fmt.Errorf("ошибка при получении оборудования: %v", err)
	}
	log.Printf("ID: %d, Model: %s, Status: %s, AssignedTo: %v", equipment.ID, equipment.Model, equipment.Status, equipment.AssignedTo)

	return &equipment, nil
}

// CreateEquipment создает новую единицу оборудования в базе данных
func (s *EquipmentService) CreateEquipment(model, serialNumber, status string) (*Equipment, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO equipment (model, serial_number, status) VALUES ($1, $2, $3) RETURNING id",
		model, serialNumber, status,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании оборудования: %v", err)
	}

	return &Equipment{ID: id, Model: model, SerialNumber: serialNumber, Status: status}, nil
}

// GetEquipmentByID возвращает оборудование по его идентификатору
func (s *EquipmentService) GetEquipmentByID(id int) (*Equipment, error) {
	var equipment Equipment
	err := s.db.QueryRow(
		"SELECT id, model, serial_number, status, assigned_to FROM equipment WHERE id = $1", id,
	).Scan(&equipment.ID, &equipment.Model, &equipment.SerialNumber, &equipment.Status, &equipment.AssignedTo)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("оборудование с id %d не найдено", id)
		}
		return nil, fmt.Errorf("ошибка при получении оборудования: %v", err)
	}

	return &equipment, nil
}

// GetAllEquipment возвращает список всех единиц оборудования
func (s *EquipmentService) GetAllEquipment() ([]Equipment, error) {
	rows, err := s.db.Query("SELECT id, model, serial_number, status, assigned_to FROM equipment")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех единиц оборудования: %v", err)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			fmt.Printf("Ошибка при закрытии rows: %v\n", cerr)
		}
	}()

	var equipmentList []Equipment
	for rows.Next() {
		var equipment Equipment
		if err := rows.Scan(&equipment.ID, &equipment.Model, &equipment.SerialNumber, &equipment.Status, &equipment.AssignedTo); err != nil {
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
func (s *EquipmentService) UpdateEquipment(id int, model, serialNumber, status string) error {
	_, err := s.db.Exec(
		"UPDATE equipment SET model = $1, serial_number = $2, status = $3 WHERE id = $4",
		model, serialNumber, status, id,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении оборудования: %v", err)
	}

	return nil
}

// DeleteEquipment удаляет оборудование из базы данных
func (s *EquipmentService) DeleteEquipment(id int) error {
	_, err := s.db.Exec("DELETE FROM equipment WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении оборудования: %v", err)
	}

	return nil
}

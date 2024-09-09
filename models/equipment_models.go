package models

// Equipment представляет сущность оборудования
type Equipment struct {
	ID           int    `json:"id"`
	Model        string `json:"model"`
	SerialNumber string `json:"serial_number"`
	Status       string `json:"status"`
	AssignedTo   *int   `json:"assigned_to"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// EquipmentRepository описывает интерфейс для работы с оборудованием
type EquipmentRepository interface {
	CreateEquipment(equipment *Equipment) (*Equipment, error)
	GetEquipmentByID(id int) (*Equipment, error)
	GetAllEquipment() ([]Equipment, error)
	UpdateEquipment(equipment *Equipment) error
	DeleteEquipment(id int) error
}

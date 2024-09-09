package mocks

import (
	"inva/models"

	"github.com/stretchr/testify/mock"
)

// MockEquipmentRepository - мок для EquipmentRepository
type MockEquipmentRepository struct {
	mock.Mock
}

// CreateEquipment создает новое оборудование
func (m *MockEquipmentRepository) CreateEquipment(equipment *models.Equipment) (*models.Equipment, error) {
	args := m.Called(equipment)
	return args.Get(0).(*models.Equipment), args.Error(1)
}

// GetEquipmentByID получает оборудование по ID
func (m *MockEquipmentRepository) GetEquipmentByID(id int) (*models.Equipment, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Equipment), args.Error(1)
}

// GetAllEquipment получает список всего оборудования
func (m *MockEquipmentRepository) GetAllEquipment() ([]models.Equipment, error) {
	args := m.Called()
	return args.Get(0).([]models.Equipment), args.Error(1)
}

// UpdateEquipment обновляет данные оборудования по ID
func (m *MockEquipmentRepository) UpdateEquipment(id int, model, status string) error {
	args := m.Called(id, model, status)
	return args.Error(0)
}

// DeleteEquipment удаляет оборудование по ID
func (m *MockEquipmentRepository) DeleteEquipment(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

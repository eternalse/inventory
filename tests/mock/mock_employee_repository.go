package mocks

import (
	"inva/models"

	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepository -  мок для EmployeeRepository
type MockEmployeeRepository struct {
	mock.Mock
}

// CreateEmployee создает нового сотрудника
func (m *MockEmployeeRepository) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	args := m.Called(employee)
	return args.Get(0).(*models.Employee), args.Error(1)
}

// GetEmployeeByID получает сотрудника по ID
func (m *MockEmployeeRepository) GetEmployeeByID(id int) (*models.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Employee), args.Error(1)
}

// GetAllEmployees получает список всех сотрудников
func (m *MockEmployeeRepository) GetAllEmployees() ([]models.Employee, error) {
	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)
}

// UpdateEmployee обновляет данные сотрудника по ID
func (m *MockEmployeeRepository) UpdateEmployee(id int, name string) error {
	args := m.Called(id, name)
	return args.Error(0)
}

// DeleteEmployee удаляет сотрудника по ID
func (m *MockEmployeeRepository) DeleteEmployee(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

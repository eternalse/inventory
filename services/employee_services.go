package services

import (
	"database/sql"
	"fmt"
	"inva/models"
	"log"
)

// Employee представляет модель сотрудника
type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Интерфейс для EmployeeRepository
type EmployeeRepository interface {
	CreateEmployee(employee *models.Employee) (*models.Employee, error)
	GetEmployeeByID(id int) (*models.Employee, error)
	GetAllEmployees() ([]models.Employee, error)
	UpdateEmployee(id int, name string) error
	DeleteEmployee(id int) error
}

// EmployeeService предоставляет методы для работы с сотрудниками
type EmployeeService struct {
	db *sql.DB
}

// NewEmployeeService создает новый экземпляр EmployeeService
func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

// CreateEmployee создает нового сотрудника в базе данных
func (s *EmployeeService) CreateEmployee(employee *Employee) (*Employee, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO employees (name) VALUES ($1) RETURNING id",
		employee.Name,
	).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при создании сотрудника: %v", err)
		return nil, err
	}

	employee.ID = id
	return employee, nil
}

// GetEmployeeByID возвращает сотрудника по его идентификатору
func (s *EmployeeService) GetEmployeeByID(id int) (*Employee, error) {
	var employee Employee
	err := s.db.QueryRow(
		"SELECT id, name FROM employees WHERE id = $1",
		id,
	).Scan(&employee.ID, &employee.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("сотрудник с id %d не найден", id)
		}
		log.Printf("Ошибка при получении сотрудника: %v", err)
		return nil, err
	}

	return &employee, nil
}

// GetAllEmployees возвращает всех сотрудников из базы данных
func (s *EmployeeService) GetAllEmployees() ([]Employee, error) {
	rows, err := s.db.Query("SELECT id, name FROM employees")
	if err != nil {
		log.Printf("Ошибка при получении всех сотрудников: %v", err)
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.Name); err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return nil, err
		}
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при переборе строк: %v", err)
		return nil, err
	}

	return employees, nil
}

// UpdateEmployee обновляет данные сотрудника
func (s *EmployeeService) UpdateEmployee(id int, name string) error {
	_, err := s.db.Exec(
		"UPDATE employees SET name = $1 WHERE id = $2",
		name,
		id,
	)
	if err != nil {
		log.Printf("Ошибка при обновлении сотрудника: %v", err)
		return err
	}

	return nil
}

// DeleteEmployee удаляет сотрудника из базы данных
func (s *EmployeeService) DeleteEmployee(id int) error {
	_, err := s.db.Exec(
		"DELETE FROM employees WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("Ошибка при удалении сотрудника: %v", err)
		return err
	}

	return nil
}

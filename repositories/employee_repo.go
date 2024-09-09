package repositories

import (
	"database/sql"
	"inva/models"
)

// EmployeeRepository интерфейс для репозитория сотрудников
type EmployeeRepository interface {
	CreateEmployee(employee *models.Employee) (*models.Employee, error)
	GetEmployeeByID(id int) (*models.Employee, error)
	GetAllEmployees() ([]models.Employee, error)
	UpdateEmployee(id int, name string) error
	DeleteEmployee(id int) error
}

// PostgresEmployeeRepository реализует интерфейс EmployeeRepository для работы с PostgreSQL
type PostgresEmployeeRepository struct {
	db *sql.DB
}

// NewPostgresEmployeeRepository создаёт новый экземпляр PostgresEmployeeRepository
func NewPostgresEmployeeRepository(db *sql.DB) *PostgresEmployeeRepository {
	return &PostgresEmployeeRepository{db: db}
}

// CreateEmployee создает нового сотрудника в базе данных
func (r *PostgresEmployeeRepository) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	err := r.db.QueryRow(
		"INSERT INTO employees (name, position, email, phone) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		employee.Name, employee.Position, employee.Email).Scan(&employee.ID, &employee.CreatedAt)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

// GetEmployeeByID возвращает сотрудника по его ID
func (r *PostgresEmployeeRepository) GetEmployeeByID(id int) (*models.Employee, error) {
	employee := &models.Employee{}
	err := r.db.QueryRow("SELECT id, name, position, email, phone, created_at FROM employees WHERE id = $1", id).
		Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Email, &employee.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return employee, nil
}

// GetAllEmployees возвращает список всех сотрудников
func (r *PostgresEmployeeRepository) GetAllEmployees() ([]models.Employee, error) {
	rows, err := r.db.Query("SELECT id, name, position, email, phone, created_at FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Email, &employee.CreatedAt); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// UpdateEmployee обновляет информацию о сотруднике
func (r *PostgresEmployeeRepository) UpdateEmployee(employee *models.Employee) error {
	_, err := r.db.Exec(
		"UPDATE employees SET name = $1, position = $2, email = $3, phone = $4 WHERE id = $5",
		employee.Name, employee.Position, employee.Email, employee.ID)

	return err
}

// DeleteEmployee удаляет сотрудника по его ID
func (r *PostgresEmployeeRepository) DeleteEmployee(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id = $1", id)
	return err
}

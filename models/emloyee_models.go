package models

// Employee представляет сущность сотрудника
type Employee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Position  string `json:"position"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// EmployeeRepository описывает интерфейс для работы с сотрудниками
type EmployeeRepository interface {
	CreateEmployee(employee *Employee) (*Employee, error)
	GetEmployeeByID(id int) (*Employee, error)
	GetAllEmployees() ([]Employee, error)
	UpdateEmployee(employee *Employee) error
	DeleteEmployee(id int) error
}

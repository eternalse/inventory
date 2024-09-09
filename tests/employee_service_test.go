// tests/employee_service_test.go
package services_test

import (
	"inva/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEmployeeService(db)

	// Определяем ожидаемые данные
	newEmployee := &services.Employee{Name: "John Doe"}
	mock.ExpectQuery("INSERT INTO employees").
		WithArgs(newEmployee.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызываем метод
	result, err := service.CreateEmployee(newEmployee)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, newEmployee.Name, result.Name)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestGetEmployeeByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEmployeeService(db)

	// Определяем ожидаемые данные
	employee := &services.Employee{ID: 1, Name: "John Doe"}
	mock.ExpectQuery("SELECT id, name FROM employees WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(employee.ID, employee.Name))

	// Вызываем метод
	result, err := service.GetEmployeeByID(1)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, employee, result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestUpdateEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEmployeeService(db)

	// Определяем ожидаемые данные
	mock.ExpectExec("UPDATE employees SET name = \\$1 WHERE id = \\$2").
		WithArgs("John Smith", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем метод
	err = service.UpdateEmployee(1, "John Smith")

	// Проверяем результаты
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestDeleteEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEmployeeService(db)

	// Определяем ожидаемые данные
	mock.ExpectExec("DELETE FROM employees WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем метод
	err = service.DeleteEmployee(1)

	// Проверяем результаты
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestGetAllEmployees(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEmployeeService(db)

	// Определяем ожидаемые данные
	employeeList := []services.Employee{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, emp := range employeeList {
		rows.AddRow(emp.ID, emp.Name)
	}
	mock.ExpectQuery("SELECT id, name FROM employees").WillReturnRows(rows)

	// Вызываем метод
	result, err := service.GetAllEmployees()

	// Проверяем результаты
	assert.NoError(t, err)
	assert.ElementsMatch(t, employeeList, result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

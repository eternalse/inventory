package services_test

import (
	"inva/models"
	"inva/services"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateEquipment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEquipmentService(db)

	// Определяем ожидаемые данные
	newEquipment := &models.Equipment{ID: 1, Model: "Laptop", SerialNumber: "1234", Status: "available"}
	mock.ExpectQuery("INSERT INTO equipment").
		WithArgs(newEquipment.Model, newEquipment.SerialNumber, newEquipment.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Вызываем метод
	result, err := service.CreateEquipment(newEquipment.Model, newEquipment.SerialNumber, newEquipment.Status)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, newEquipment.ID, result.ID)
	assert.Equal(t, newEquipment.Model, result.Model)
	assert.Equal(t, newEquipment.SerialNumber, result.SerialNumber)
	assert.Equal(t, newEquipment.Status, result.Status)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestGetEquipmentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEquipmentService(db)

	// Ожидаемые данные
	rows := sqlmock.NewRows([]string{"id", "model", "serial_number", "status", "assigned_to"}).
		AddRow(1, "Laptop", "1234", "available", nil)

	// Ожидаемый запрос
	mock.ExpectQuery("^SELECT id, model, serial_number, status, assigned_to FROM equipment WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(rows)

	// Вызываем метод
	result, err := service.GetEquipmentByID(1)

	// Проверяем результаты
	expected := &services.Equipment{
		ID:           1,
		Model:        "Laptop",
		SerialNumber: "1234",
		Status:       "available",
		AssignedTo:   nil,
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestUpdateEquipment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEquipmentService(db)

	// Определяем ожидаемые данные
	mock.ExpectExec(`^UPDATE equipment SET model = \$1, serial_number = \$2, status = \$3 WHERE id = \$4$`).
		WithArgs("Laptop", "1234", "available", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем метод
	err = service.UpdateEquipment(1, "Laptop", "1234", "available")

	// Проверяем результаты
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestDeleteEquipment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEquipmentService(db)

	// Определяем ожидаемые данные
	mock.ExpectExec("DELETE FROM equipment WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Вызываем метод
	err = service.DeleteEquipment(1)

	// Проверяем результаты
	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

func TestGetAllEquipment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	service := services.NewEquipmentService(db)

	// Определяем ожидаемые данные
	equipmentList := []services.Equipment{
		{ID: 1, Model: "Laptop", SerialNumber: "1234", Status: "available"},
		{ID: 2, Model: "Monitor", SerialNumber: "5678", Status: "in use"},
	}
	rows := sqlmock.NewRows([]string{"id", "model", "serial_number", "status", "assigned_to"})
	for _, eq := range equipmentList {
		rows.AddRow(eq.ID, eq.Model, eq.SerialNumber, eq.Status, nil)
	}
	mock.ExpectQuery("SELECT id, model, serial_number, status, assigned_to FROM equipment").WillReturnRows(rows)

	// Вызываем метод
	result, err := service.GetAllEquipment()

	// Проверяем результаты
	assert.NoError(t, err)
	assert.ElementsMatch(t, equipmentList, result)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Ожидания не были удовлетворены: %v", err)
	}
}

# Equipment Management System

This project is a simple system for managing equipment and assigning it to users. It allows you to create, update, and fetch information about equipment stored in a PostgreSQL database. The project also supports assigning equipment to users and tracking who is responsible for specific items. It is written in Go and includes unit tests with `sqlmock` for database interaction testing.

**Note:** This project is still in development. Features may change, and new functionality may be added. Please check back for updates!

## Features

- Create new equipment records.
- Update existing equipment information.
- Fetch equipment details by ID.
- Assign equipment to users.
- Track which user is assigned to specific equipment.
- Unit tests for key functionalities.

## Technologies

- **Go** - Core programming language.
- **PostgreSQL** - Database used for storing equipment and user information.
- **sqlmock** - Used for mocking database interactions in unit tests.

## Installation

### . Clone the repository

```bash
git clone https://github.com/eternalse/inventory.git
cd inventory

## Database Schema

Make sure to set up the following tables in your PostgreSQL database:

### 1. Equipment Table

| Column        | Type             | Description                           |
|---------------|------------------|---------------------------------------|
| id            | INTEGER          | Primary Key, Auto-increment           |
| model         | TEXT             | Equipment model name                  |
| status        | CHARACTER VARYING| Status of the equipment               |
| assigned_to   | INTEGER          | (Optional) Foreign key to users table |

### Relationships

- The `assigned_to` column in the `equipment` table is a foreign key that references the `id` column in the `users` table.


### 2. Employees Table 

| Column     | Type         | Description                    |
|------------|--------------|--------------------------------|
| id         | INTEGER      | Primary Key, Auto-increment    |
| name       | TEXT         | Users name                     |


### 3. Serial Numbers Table

| Column         | Type             | Description                                           |
|----------------|------------------|-------------------------------------------------------|
| id             | SERIAL           | Primary Key, Auto-increment                           |
| serial_number  | VARCHAR(100)     | Unique serial number, must be unique                  |
| equipment_id   | INT              | Foreign key referencing the `id` in `equipment` table |


### 4. Equipment Logs Table

| Column        | Type              | Description                                                            |
|---------------|-------------------|------------------------------------------------------------------------|
| id            | SERIAL            | Primary Key, Auto-increment                                            |
| equipment_id  | INT               | Foreign key referencing the `id` in the `equipment` table              |
| user_id       | INT               | Foreign key referencing the `id` in the `employees` table              |
| issued_at     | TIMESTAMP         | Timestamp when the equipment was issued, defaults to current timestamp |
| returned_at   | TIMESTAMP         | Timestamp when the equipment was returned (nullable)                   |
| status        | VARCHAR(50)       | Status of the equipment (e.g., 'issued', 'returned')                   |

### Description:
id — Auto-incrementing primary key.
equipment_id — Foreign key referencing the equipment table, indicating which piece of equipment was issued.
user_id — Foreign key referencing the employees table, indicating to which user the equipment was issued.
issued_at — Timestamp when the equipment was issued, defaults to the current timestamp.
returned_at — Timestamp when the equipment was returned. This column can be nullable.
status — String indicating the status of the equipment (e.g., 'issued' or 'returned').


## Example Commands

1. Creating a new employee

     curl -X POST http://localhost:8080/employees \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe"}'


2. Getting a list of all employees

     curl http://localhost:8080/employees


3. Getting an employee by ID

     curl http://localhost:8080/employees/1

4. Updating employee information

    curl -X PUT http://localhost:8080/employees/9 \
       -H "Content-Type: application/json" \
       -d '{"name": "Jane Doe"}'

5. Deleting an employee by ID

    curl -X DELETE http://localhost:8080/employees/9
    
 ## Managing Equipment
1. Creating new equipment   
    
    curl -X POST http://localhost:8080/equipment \
     -H "Content-Type: application/json" \
     -d '{
           "model": "Laptop",
           "status": "available",
           "serial_number": "ABC123",
           "description": "A high-end laptop"
         }'

2. Getting all equipment records

   curl -X GET http://localhost:8080/equipment


3. Getting equipment by ID
 
   curl -X GET http://localhost:8080/equipment/1

4. Updating equipment information
 
    curl -X PUT http://localhost:8080/equipment/12 \
     -H "Content-Type: application/json" \
     -d '{
           "model": "Laptop Pro",
           "status": "in use",
           "serial_number": "ABC1234",
           "description": "An updated high-end laptop"
         }'

5. Deleting equipment

   curl -X DELETE http://localhost:8080/equipment/1


6. Assigning equipment to a user

   curl -X POST http://localhost:8080/equipment/7/assign/user/5


7. Getting details of equipment assigned to a user

   curl -X GET http://localhost:8080/equipment/7/details


8. Returning equipment from a user

   curl -X PUT http://localhost:8080/equipment/6/return



### Description of Parameters
ID — Unique identifier for employees or equipment.
Model — Model of the equipment.
Status — Current status of the equipment (e.g., "available", "in use").
Serial Number — Serial number of the equipment.
Description — Description of the equipment.
User ID — Unique identifier for the use






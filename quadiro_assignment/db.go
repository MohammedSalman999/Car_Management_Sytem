package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	// Replace with your MySQL connection details
	db, err = sql.Open("mysql", "localhost:Zeeshan1khan$@tcp(127.0.0.1:3306)/quadiro_db")
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	createTables()
}

func createTables() {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS cars (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			manufacturing_year INT NOT NULL,
			price DECIMAL(10, 2) NOT NULL
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS admins (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL
		);
		`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating tables: %v\n", err)
		}
	}
}

func getAllCars() ([]Car, error) {
	rows, err := db.Query("SELECT id, name, manufacturing_year, price FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		err := rows.Scan(&car.ID, &car.Name, &car.ManufacturingYear, &car.Price)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	defaultDBHost = "localhost"
	defaultDBPort = "3306"
	defaultDBUser = "root"
	defaultDBPass = ""
	defaultDBName = "cetec"
)

func GetDBConnection() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load .env")
	}
	dbHost := getEnv("DB_HOST", defaultDBHost)
	dbPort := getEnv("DB_PORT", defaultDBPort)

	cfg := mysql.Config{
		User:                 getEnv("DB_USER", defaultDBUser),
		Passwd:               getEnv("DB_PASS", defaultDBPass),
		DBName:               getEnv("DB_NAME", defaultDBName),
		Addr:                 dbHost + ":" + dbPort,
		AllowNativePasswords: true,
		ParseTime:            true,
		Net:                  "tcp",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")
	return db, nil
}

func GetPersonInfo(db *sql.DB, personID int) (*PersonInfo, error) {
	query := `
		SELECT p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code
		FROM person p
		JOIN phone ph ON p.id = ph.person_id
		JOIN address_join aj ON p.id = aj.person_id
		JOIN address a ON aj.address_id = a.id
		WHERE p.id = ?
	`

	var personInfo PersonInfo
	err := db.QueryRow(query, personID).Scan(
		&personInfo.Name,
		&personInfo.PhoneNumber,
		&personInfo.City,
		&personInfo.State,
		&personInfo.Street1,
		&personInfo.Street2,
		&personInfo.ZipCode,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("person with ID %d not found", personID)
		}
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	return &personInfo, nil
}

func CreatePerson(db *sql.DB, person *PersonCreate) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	var personID int64
	result, err := tx.Exec("INSERT INTO person (name, age) VALUES (?, ?)", person.Name, 0)
	if err != nil {
		return fmt.Errorf("failed to insert person: %v", err)
	}
	personID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID for person: %v", err)
	}

	_, err = tx.Exec("INSERT INTO phone (number, person_id) VALUES (?, ?)", person.PhoneNumber, personID)
	if err != nil {
		return fmt.Errorf("failed to insert phone: %v", err)
	}

	result, err = tx.Exec(
		"INSERT INTO address (city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)",
		person.City, person.State, person.Street1, person.Street2, person.ZipCode,
	)
	if err != nil {
		return fmt.Errorf("failed to insert address: %v", err)
	}
	addressID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID for address: %v", err)
	}

	_, err = tx.Exec("INSERT INTO address_join (person_id, address_id) VALUES (?, ?)", personID, addressID)
	if err != nil {
		return fmt.Errorf("failed to insert address join: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
	UpdateAccount(*Account) error
	DeleteAccount(int) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// Connect to the database
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	// Create vinyl_store_storage database
	if err := s.CreateDatabase(); err != nil {
		return err
	}

	// Create vinyl_user role
	if err := s.CreateUser(); err != nil {
		return err
	}

	// Create vinyl_store schema
	if err := s.CreateSchema(); err != nil {
		return err
	}

	// Grant permissions to vinyl_user
	if err := s.GrantPermissions(); err != nil {
		return err
	}

	// Connect to the database with the new user
	s.db.Close() // Close the current connection so we can reconnect with the new user.

	connStr := "user=vinyl_user dbname=vinyl_store_storage password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Assign the new connection to the store
	s.db = db

	// Create accounts table
	if err := s.CreateAccountTable(); err != nil {
		return err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := "INSERT INTO vinyl_store.accounts (first_name, last_name, email, created_at) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, a.FirstName, a.LastName, a.Email, a.CreatedAt)
	return err
}
func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM vinyl_store.accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// Helper functions for DB setup
func (s *PostgresStore) CreateDatabase() error {
	_, err := s.db.Exec("CREATE DATABASE vinyl_store_storage")
	if err != nil && err.Error() != "pq: database \"vinyl_store_storage\" already exists" {
		return err
	}
	// Close the current connection
	s.db.Close()

	// Open a new connection to the vinyl_store_storage database
	connStr := "user=postgres dbname=vinyl_store_storage password=mysecretpassword sslmode=disable"
	s.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test the new connection
	err = s.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateUser() error {
	_, err := s.db.Exec("CREATE ROLE vinyl_user WITH LOGIN PASSWORD 'mysecretpassword'")
	if err != nil && err.Error() != "pq: role \"vinyl_user\" already exists" {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateSchema() error {
	query := `
        CREATE SCHEMA IF NOT EXISTS vinyl_store;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GrantPermissions() error {
	query := `
        GRANT CONNECT ON DATABASE vinyl_store_storage TO vinyl_user;
        GRANT CREATE ON DATABASE vinyl_store_storage TO vinyl_user;
        GRANT USAGE ON SCHEMA vinyl_store TO vinyl_user;
        GRANT CREATE ON SCHEMA vinyl_store TO vinyl_user;
        GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA vinyl_store TO vinyl_user;
        ALTER DEFAULT PRIVILEGES IN SCHEMA vinyl_store GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO vinyl_user;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS vinyl_store.accounts (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			email VARCHAR(255),
			created_at TIMESTAMP DEFAULT NOW()
		)
	`
	_, err := s.db.Exec(query)
	return err
}

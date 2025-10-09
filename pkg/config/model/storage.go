package model

import "fmt"

// Storage Configuration
type Storage struct {
	PostgreSQL Database
	Redis      Redis
}

// Database Configuration
type Database struct {
	Address  string
	Port     int
	Database string
	Username string
	Password string
	SSLMode  bool
	TimeZone string
}

type DatabaseDriver string

const PostgreSQLDriver DatabaseDriver = "postgres"

// GetDSN get DSN by driver database
func (d Database) GetDSN(driver DatabaseDriver) string {
	if driver == PostgreSQLDriver {
		sslMode := "disable"
		if d.SSLMode {
			sslMode = "enable"
		}
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			d.Address, d.Username, d.Password, d.Database, d.Port, sslMode, d.TimeZone)
	}
	return ""
}

func (d Database) GetDatabaseURL(driver DatabaseDriver) string {
	if driver == PostgreSQLDriver {
		sslMode := "disable"
		if d.SSLMode {
			sslMode = "enable"
		}
		return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
			driver, d.Username, d.Password, d.Address, d.Port, d.Database, sslMode)
	}
	return ""
}

// Redis Configuration
type Redis struct {
	Address string
	Port    int
}

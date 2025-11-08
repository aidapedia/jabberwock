package model

// Storage Configuration
type Storage struct {
	Redis Redis
}

// Redis Configuration
type Redis struct {
	Address string
	Port    int
}

package model

// SecretConfig is the configuration for the secret
type SecretConfig struct {
	Database struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
	} `json:"database"`
}

// ServiceConfig is the configuration for the service
type ServiceConfig struct {
	App     App
	Storage Storage
	Vendor  Vendor
}

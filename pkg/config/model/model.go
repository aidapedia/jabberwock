package model

// SecretConfig is the configuration for the secret
type SecretConfig struct {
	Auth struct {
		PrivateKey string `json:"private_key"`
	} `json:"auth"`
	Database struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
	} `json:"database"`
	Tracer struct {
		AddressURL string `json:"address_url"`
	} `json:"tracer"`
}

// ServiceConfig is the configuration for the service
type ServiceConfig struct {
	App     App
	Storage Storage
	Vendor  Vendor
	// Secret Config
	Secret SecretConfig
}

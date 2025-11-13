package model

// SecretConfig is the configuration for the secret
type SecretConfig struct {
	Auth struct {
		PrivateKey string
		PublicKey  string
	}
	Database struct {
		Name     string
		Username string
		Password string
		Host     string
		Port     int
	}
	Tracer struct {
		AddressURL string
	}
}

// ServiceConfig is the configuration for the service
type ServiceConfig struct {
	App     App
	Storage Storage
	Vendor  Vendor
	// Secret Config
	Secret SecretConfig
}

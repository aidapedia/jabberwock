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
	Redis struct {
		Address string
		Port    int
	}
	Tracer struct {
		AddressURL string
	}
	Whatsapp struct {
		Host    string
		APIKey  string
		Session string
	}
}

// ServiceConfig is the configuration for the service
type ServiceConfig struct {
	// Application Config
	App App
	// Secret Config
	Secret SecretConfig
	// Feature Flags
	FeatureFlags struct {
		DisableTracer bool
	}
}

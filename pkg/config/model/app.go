package model

// App Configuration
type App struct {
	// Name of the application should be same like your github repository name
	// For example: github.com/aidapedia/jabberwock
	// The name should be "jabberwock"
	Name string
	// LocalTime is the time zone of the application
	LocalTime string
	Log       struct {
		// Level of the logger
		Level string
	}
	// HTTPServer Configuration
	HTTPServer struct {
		// Address of the HTTP server
		Address string
		// Port of the HTTP server
		Port int
	}

	Auth struct {
		Issuer     string
		ModelPath  string
		PolicyPath string
	}
}

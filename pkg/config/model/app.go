package model

// App Configuration
type App struct {
	Name       string
	HTTPServer struct {
		Address string
		Port    int
	}

	Auth struct {
		ModelPath  string
		PolicyPath string
	}
}

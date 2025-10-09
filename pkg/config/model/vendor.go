package model

// Vendor Configuration
type Vendor struct {
	Midtrans Midtrans
	Twilio   Twilio
	SendGrid SendGrid
}

// Midtrans Configuration
type Midtrans struct {
	IsProduction bool
	ServerKey    string
}

// Twilio Configuration
type Twilio struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
}

// SendGrid Configuration
type SendGrid struct {
	APIKey string
}

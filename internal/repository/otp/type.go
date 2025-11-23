package otp

import "time"

type RegistrationOTP struct {
	RegistrationData RegistrationData
	OTP              string
	RequestAt        time.Time
}

type RegistrationData struct {
	Phone    string
	Name     string
	Password string
}

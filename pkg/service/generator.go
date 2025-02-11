package service

import (
	"github.com/Aryagorjipour/SPWD/pkg/domain"
)

// ValidateMode converts a string mode into a Mode type and validates it.
func ValidateMode(m string) (domain.Mode, error) {
	mode, err := domain.NewMode(m)
	if err != nil {
		return "", err
	}
	return mode, nil
}

// ValidateLength checks if the password length is valid for the given mode.
func ValidateLength(length int, mode domain.Mode) error {
	return domain.ValidateLength(length, mode)
}

// GeneratePassword generates the password based on the length and mode.
func GeneratePassword(length int, mode domain.Mode) (string, error) {
	// Validate the length based on the mode
	err := ValidateLength(length, mode)
	if err != nil {
		return "", err
	}

	// Generate and return the password
	password := domain.GeneratePassword(length, mode)
	return password, nil
}

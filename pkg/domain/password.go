package domain

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

// Mode defines different password generation modes.
type Mode string

// Constants representing different password modes with the short flags
const (
	VeryWeak    Mode = "vw" // very weak
	Weak        Mode = "w"  // weak
	Medium      Mode = "m"  // medium
	Strong      Mode = "s"  // strong
	VeryStrong  Mode = "vs" // very strong
	Unbreakable Mode = "xb" // unbreakable
)

// NewMode validates the mode and returns it based on the flag.
func NewMode(m string) (Mode, error) {
	switch m {
	case "vw":
		return VeryWeak, nil
	case "w":
		return Weak, nil
	case "m":
		return Medium, nil
	case "s":
		return Strong, nil
	case "vs":
		return VeryStrong, nil
	case "xb":
		return Unbreakable, nil
	default:
		return "", errors.New("invalid mode")
	}
}

// GeneratePassword generates a password based on the mode and length.
func GeneratePassword(length int, mode Mode) string {
	rand.Seed(time.Now().UnixNano())

	var password strings.Builder
	var charset string

	// Define character sets for each mode
	switch mode {
	case VeryWeak:
		// Very Weak - only digits (easy to guess)
		charset = "0123456789"
	case Weak:
		// Weak - digits and lowercase letters
		charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	case Medium:
		// Medium - digits, lowercase, and uppercase letters
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case Strong:
		// Strong - digits, lowercase, uppercase, and special characters
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=<>?"
	case VeryStrong:
		// Very Strong - digits, lowercase, uppercase, and special characters with added complexity
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=<>?{}[]:;,.<>"
	case Unbreakable:
		// Unbreakable - includes every character type, higher randomness, and length
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_-+=<>?{}[]:;,.<>~`|"
	}

	// Ensure that password length is valid for the selected mode
	if length < 1 {
		return ""
	}

	// Generate the password by picking random characters from the charset
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		password.WriteByte(charset[randomIndex])
	}

	// Return the generated password
	return password.String()
}

// ValidateLength validates if the length is within the acceptable range for the given mode.
func ValidateLength(length int, mode Mode) error {
	switch mode {
	case VeryWeak:
		if length < 4 || length > 6 {
			return errors.New("length for very weak password must be between 4 and 6")
		}
	case Weak:
		if length < 6 || length > 8 {
			return errors.New("length for weak password must be between 6 and 8")
		}
	case Medium:
		if length < 8 || length > 12 {
			return errors.New("length for medium password must be between 8 and 12")
		}
	case Strong:
		if length < 12 || length > 16 {
			return errors.New("length for strong password must be between 12 and 16")
		}
	case VeryStrong:
		if length < 16 || length > 20 {
			return errors.New("length for very strong password must be between 16 and 20")
		}
	case Unbreakable:
		if length < 20 || length > 30 {
			return errors.New("length for unbreakable password must be between 20 and 30")
		}
	}
	return nil
}

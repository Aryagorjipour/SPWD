package storage

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/nacl/secretbox"
	"os"
	"runtime"
)

type Config struct {
	SecretKey string `json:"secret_key"`
}

var secretKey [32]byte

// GetConfigPath determines the config path based on OS
func GetConfigPath() string {
	if runtime.GOOS == "windows" {
		return "C:\\ProgramData\\spwd\\config.json"
	}
	return "/etc/spwd/config.json"
}

// LoadSecretKey loads the secret key from config.json
func LoadSecretKey() error {
	configPath := GetConfigPath()

	file, err := os.Open(configPath)
	if err != nil {
		return errors.New("failed to open config.json: " + err.Error())
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return errors.New("failed to parse config.json: " + err.Error())
	}

	if len(config.SecretKey) != 32 {
		return errors.New("invalid secret key length: must be 32 bytes")
	}

	copy(secretKey[:], config.SecretKey)
	return nil
}

// Encrypt encrypts the password
func Encrypt(plaintext string) (string, error) {

	//Ensure the secret key is loaded before encryption
	if err := LoadSecretKey(); err != nil {
		return "", err
	}

	var nonce [24]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return "", err
	}

	encrypted := secretbox.Seal(nonce[:], []byte(plaintext), &nonce, &secretKey)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts the password
func Decrypt(ciphertext string) (string, error) {
	// Ensure the secret key is loaded before decryption
	if err := LoadSecretKey(); err != nil {
		return "", err
	}

	if len(ciphertext) == 0 {
		return "", errors.New("empty encrypted data")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", errors.New("failed to decode encrypted data")
	}

	if len(data) < 24 { // Ensure it contains nonce + encrypted text
		return "", errors.New("invalid encrypted data length")
	}

	var nonce [24]byte
	copy(nonce[:], data[:24])

	decrypted, ok := secretbox.Open(nil, data[24:], &nonce, &secretKey)
	if !ok {
		return "", errors.New("decryption failed")
	}

	return string(decrypted), nil
}

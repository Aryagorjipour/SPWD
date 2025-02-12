package storage

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
	"os"
	"path/filepath"
)

// Config structure for loading the secret key
type Config struct {
	SecretKey string `json:"secret_key"`
}

var secretKey [32]byte

// GetConfigPath determines the correct config.json path based on the executable location
func GetConfigPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Join(filepath.Dir(exePath), "config.json")
}

// LoadSecretKey loads the key from config.json and decodes it
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

	// Decode base64 secret key
	decodedKey, err := base64.StdEncoding.DecodeString(config.SecretKey)
	if err != nil {
		return errors.New("failed to decode base64 secret key: " + err.Error())
	}

	// Ensure the secret key is exactly 32 bytes
	if len(decodedKey) != 32 {
		return errors.New(fmt.Sprintf("invalid secret key length: %d; expected 32", len(decodedKey)))
	}

	copy(secretKey[:], decodedKey)
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

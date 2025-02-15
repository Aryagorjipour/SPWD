package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

// GetExecutablePath returns the directory where the binary is located
func GetExecutablePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// OpenDB initializes the database in the same directory as the executable
func OpenDB() error {
	exeDir, err := GetExecutablePath()
	if err != nil {
		fmt.Println("❌ Error determining executable path:", err)
		return err
	}

	dbPath := filepath.Join(exeDir, "passwords.db")

	// Open the database
	var errDB error
	db, errDB = bbolt.Open(dbPath, 0660, nil)
	if errDB != nil {
		fmt.Println("❌ Failed to open database:", errDB)
		return errDB
	}

	// Create "Passwords" bucket
	errBucket := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Passwords"))
		return err
	})

	if errBucket != nil {
		fmt.Println("❌ Failed to create Passwords bucket:", errBucket)
		return errBucket
	}

	return nil
}

var db *bbolt.DB

// PasswordEntry represents a stored password
type PasswordEntry struct {
	ID        int    `json:"id"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	Note      string `json:"note,omitempty"`
}

// SavePassword stores a password in the database with an auto-incremented ID
func SavePassword(password string) (int, error) {
	err := OpenDB()
	if err != nil {
		fmt.Println("❌ Error opening database:", err)
		return 0, err
	}
	defer db.Close()

	// Encrypt password
	encryptedPass, err := Encrypt(password)
	if err != nil {
		fmt.Println("❌ Error encrypting password:", err)
		return 0, err
	}

	var id uint64
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))
		if b == nil {
			return errors.New("passwords bucket missing")
		}

		// Get next sequence ID
		nextID, err := b.NextSequence()
		if err != nil {
			return err
		}
		id = nextID

		// Store password
		entry := PasswordEntry{
			ID:        int(id),
			Password:  encryptedPass,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		data, _ := json.Marshal(entry)

		err = b.Put([]byte(fmt.Sprintf("%d", id)), data)
		if err != nil {
			fmt.Println("❌ Failed to store password:", err)
		}
		return err
	})

	return int(id), err
}

// GetAllPasswords retrieves all passwords (decrypted)
func GetAllPasswords() ([]PasswordEntry, error) {
	err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var passwords []PasswordEntry
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))
		b.ForEach(func(k, v []byte) error {
			keyStr := string(k)

			// Skip the "lastID" entry (it's not a password)
			if keyStr == "lastID" {
				return nil
			}

			var entry PasswordEntry
			json.Unmarshal(v, &entry)

			// Ensure ID is valid (should never be 0)
			if entry.ID == 0 {
				fmt.Println("Warning: Skipping corrupted password entry with ID 0")
				return nil
			}

			// Attempt to decrypt password
			decryptedPass, err := Decrypt(entry.Password)
			if err != nil {
				fmt.Printf("Warning: Skipping corrupted password ID %d\n", entry.ID)
				return nil // Skip this entry
			}

			entry.Password = decryptedPass
			passwords = append(passwords, entry)
			return nil
		})
		return nil
	})
	return passwords, err
}

// DeletePassword removes a password by ID
func DeletePassword(id int) error {
	err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))
		return b.Delete([]byte(fmt.Sprintf("%d", id)))
	})
}

// AddNoteToPassword adds a note to a stored password
func AddNoteToPassword(id int, note string) error {
	err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))
		data := b.Get([]byte(fmt.Sprintf("%d", id)))
		if data == nil {
			return errors.New("password not found")
		}

		var entry PasswordEntry
		json.Unmarshal(data, &entry)
		entry.Note = note
		updatedData, _ := json.Marshal(entry)
		return b.Put([]byte(fmt.Sprintf("%d", id)), updatedData)
	})
}

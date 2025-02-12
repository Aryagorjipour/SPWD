package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
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
		log.Println("❌ Error determining executable path:", err)
		return err
	}

	dbPath := filepath.Join(exeDir, "passwords.db")
	log.Println("🟢 Using database path:", dbPath) // Debug log

	// Check if the database file exists, create it if necessary
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("⚠️ Database file not found. Creating new database...")
		file, err := os.Create(dbPath)
		if err != nil {
			log.Println("❌ Error creating database file:", err)
			return err
		}
		file.Close()

		// Set group read/write permissions
		err = os.Chmod(dbPath, 0660)
		if err != nil {
			log.Println("❌ Failed to set permissions on database file:", err)
		}

		// Set the correct group (same as the executable)
		exeInfo, err := os.Stat(filepath.Join(exeDir, "spwd"))
		if err == nil {
			if stat, ok := exeInfo.Sys().(*syscall.Stat_t); ok {
				os.Chown(dbPath, int(stat.Uid), int(stat.Gid))
			}
		}
	}

	// Open the database
	db, err = bbolt.Open(dbPath, 0660, nil)
	if err != nil {
		log.Println("❌ Failed to open database:", err)
		return err
	}

	// Create the "Passwords" bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Passwords"))
		return err
	})
	if err != nil {
		log.Println("❌ Failed to create Passwords bucket:", err)
		return err
	}

	log.Println("✅ Database initialized successfully at", dbPath)
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
		log.Println("❌ Error opening database:", err)
		return 0, err
	}
	defer db.Close()

	// Encrypt the password before saving
	encryptedPass, err := Encrypt(password)
	if err != nil {
		log.Println("❌ Error encrypting password:", err)
		return 0, err
	}

	var id uint64
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))

		// Get next sequence ID (auto-increment)
		nextID, err := b.NextSequence()
		if err != nil {
			log.Println("❌ Error getting next sequence ID:", err)
			return err
		}
		id = nextID

		// Create new entry
		entry := PasswordEntry{
			ID:        int(id),
			Password:  encryptedPass,
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		data, _ := json.Marshal(entry)

		// Store the password entry with its ID as the key
		err = b.Put([]byte(fmt.Sprintf("%d", id)), data)
		if err == nil {
			log.Printf("🟢 Password stored successfully with ID: %d\n", id)
		}
		return err
	})

	if err != nil {
		log.Println("❌ Failed to save password:", err)
		return 0, err
	}

	return int(id), nil
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

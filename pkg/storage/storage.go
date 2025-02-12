package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"go.etcd.io/bbolt"
)

// DB file path (OS-Specific)
var dbPath string
var db *bbolt.DB

func init() {
	if runtime.GOOS == "windows" {
		dbPath = "C:\\ProgramData\\spwd\\passwords.db" // Windows path
	} else {
		dbPath = "/etc/spwd/passwords.db" // Linux/macOS path
	}
}

// PasswordEntry represents a stored password
type PasswordEntry struct {
	ID        int    `json:"id"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	Note      string `json:"note,omitempty"`
}

// OpenDB initializes the database, creating it if it does not exist
func OpenDB() error {
	// Ensure the directory exists
	dbDir := dbPath[:len(dbPath)-len("/passwords.db")]
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err := os.MkdirAll(dbDir, 0755)
		if err != nil {
			log.Println("Error creating database directory:", err)
			return err
		}
	}

	// Open or create the database
	var err error
	db, err = bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Println("Failed to open database:", err)
		return err
	}

	// Create the "Passwords" bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Passwords"))
		return err
	})
	if err != nil {
		log.Println("Failed to create Passwords bucket:", err)
		return err
	}

	return nil
}

// SavePassword securely stores a generated password and returns a unique ID
func SavePassword(password string) (int, error) {
	err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Encrypt the password before saving
	encryptedPass, err := Encrypt(password)
	if err != nil {
		return 0, err
	}

	var id uint64
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Passwords"))

		// Get next sequence ID (auto-increment)
		nextID, err := b.NextSequence()
		if err != nil {
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
		return b.Put([]byte(fmt.Sprintf("%d", id)), data)
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

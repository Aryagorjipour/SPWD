package main

import (
	"flag"
	"fmt"
	"github.com/Aryagorjipour/SPWD/pkg/service"
	"github.com/Aryagorjipour/SPWD/pkg/storage"
	"os"
)

func main() {
	length := flag.Int("L", 12, "Password length")
	mode := flag.String("M", "m", "Password mode: vw (very weak), w (weak), m (medium), s (strong), vs (very strong), xb (unbreakable)")
	showAll := flag.Bool("s", false, "Show all stored passwords")
	deleteID := flag.Int("d", -1, "Delete password by ID")
	noteID := flag.Int("n", -1, "Password ID to add note")
	noteText := flag.String("note", "", "Note to add")
	flag.Parse()

	if *showAll {
		passwords, _ := storage.GetAllPasswords()
		fmt.Println("Stored Passwords:")
		for _, p := range passwords {
			fmt.Printf("[%d] %s - %s\n", p.ID, p.Password, p.Note)
		}
		os.Exit(0)
	}

	if *deleteID != -1 {
		storage.DeletePassword(*deleteID)
		fmt.Println("Password deleted!")
		os.Exit(0)
	}

	if *noteID != -1 {
		storage.AddNoteToPassword(*noteID, *noteText)
		fmt.Println("Note added!")
		os.Exit(0)
	}
	// Validate the mode and get the corresponding Mode type
	modeEnum, err := service.ValidateMode(*mode)
	if err != nil {
		fmt.Println("Invalid mode:", err)
		os.Exit(1)
	}

	// Validate the password length
	err = service.ValidateLength(*length, modeEnum)
	if err != nil {
		fmt.Println("Invalid length:", err)
		os.Exit(1)
	}

	// Call the generator to generate the password
	password, err := service.GeneratePassword(*length, modeEnum)
	if err != nil {
		fmt.Println("Error generating password:", err)
	}
	id, _ := storage.SavePassword(password)
	fmt.Printf("Generated Password (ID: %d): %s\n", id, password)
}

package main

import (
	"flag"
	"fmt"
	"github.com/Aryagorjipour/SPWD/pkg/service"
	"github.com/atotto/clipboard"
	"os"
)

func main() {
	// Parse flags
	length := flag.Int("L", 12, "Password length")
	mode := flag.String("M", "m", "Password mode: vw (very weak), w (weak), m (medium), s (strong), vs (very strong), xb (unbreakable)")
	flag.Parse()

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
		os.Exit(1)
	}

	// Output the password
	fmt.Println("Generated Password:", password)

	// Copy password to clipboard
	err = clipboard.WriteAll(password)
	if err != nil {
		fmt.Println("Error copying to clipboard:", err)
		os.Exit(1)
	}

	fmt.Println("Password copied to clipboard!")
}

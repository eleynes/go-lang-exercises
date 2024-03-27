package main

import (
	"flag"
	"fmt"

	"example.com/passwordgenerator"
)

func parseFlags() (int, bool, bool, bool, string, int) {
	// Define flags for user preferences
	length := flag.Int("length", 6, "Length of the password/PIN")
	includeNumbers := flag.Bool("includeNumbers", true, "Include numbers in the password")
	includeSymbols := flag.Bool("includeSymbols", true, "Include symbols in the password")
	includeUppercase := flag.Bool("includeUppercase", true, "Include uppercase letters in the password")
	passwordType := flag.String("type", "random", "Type of password (random, alphanumeric or pin)")
	numPasswords := flag.Int("count", 1, "Number of passwords to generate")

	flag.Parse()

	return *length, *includeNumbers, *includeSymbols, *includeUppercase, *passwordType, *numPasswords
}

func main() {

	length, includeNumbers, includeSymbols, includeUppercase, passwordType, numPasswords := parseFlags()

	for i := 0; i < numPasswords; i++ {
		var password string
		var err error

		if passwordType == "random" {
			password, err = passwordgenerator.GenerateSecureRandomPassword(length, includeNumbers, includeSymbols, includeUppercase)
		} else if passwordType == "pin" {
			password, err = passwordgenerator.GenerateSecurePIN(length)
		} else if passwordType == "alphanumeric" {
			password, err = passwordgenerator.GenerateSecureAlphanumericPassword(length, includeNumbers, includeUppercase)
		} else {
			fmt.Println("Invalid password type. Please choose 'random','alphanumeric' or 'pin'.")
			return
		}

		if err != nil {
			fmt.Println("Error generating password/PIN:", err)
			return
		}

		fmt.Printf("Generated Secure Password/PIN %d: %s\n", i+1, password)
	}
}

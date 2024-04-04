package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/passwordgenerator"
)

type Password struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type PasswordReq struct {
	Length              int
	NumPasswords        int
	PasswordType        string
	IsNumbersIncluded   bool
	IsSymbolsIncluded   bool
	IsUppercaseIncluded bool
}

func main() {
	http.HandleFunc("/generate-password", generatePassword)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generatePassword(w http.ResponseWriter, r *http.Request) {
	var passwords []Password
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(PasswordReq)
	err := body.Decode(&params)

	numPasswordsInt := params.NumPasswords
	passwordType := params.PasswordType
	lengthInt := params.Length
	includeNumbersBool := params.IsNumbersIncluded
	includeSymbolsBool := params.IsSymbolsIncluded
	includeUppercaseBool := params.IsUppercaseIncluded

	for i := 0; i < numPasswordsInt; i++ {
		var password string

		if passwordType == "random" {
			password, err = passwordgenerator.GenerateSecureRandomPassword(lengthInt, includeNumbersBool, includeSymbolsBool, includeUppercaseBool)
		} else if passwordType == "pin" {
			password, err = passwordgenerator.GenerateSecurePIN(lengthInt)
		} else if passwordType == "alphanumeric" {
			password, err = passwordgenerator.GenerateSecureAlphanumericPassword(lengthInt, includeNumbersBool, includeUppercaseBool)
		} else {
			fmt.Println("Invalid password type. Please choose 'random','alphanumeric' or 'pin'.")
			return
		}

		if err != nil {
			fmt.Println("Error generating password/PIN:", err)
			return
		}

		var passwordobj Password
		passwordobj.ID = i + 1
		passwordobj.Password = password

		passwords = append(passwords, passwordobj)
	}

	json.NewEncoder(w).Encode(passwords)
}

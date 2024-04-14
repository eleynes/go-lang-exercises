package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"example.com/auth"
	"example.com/db"
	"example.com/passwordgenerator"
	"golang.org/x/crypto/bcrypt"
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

type Item struct {
	VaultId  int    `json:"vaultid"`
	ItemName string `json:"itemname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
	Notes    string `json:"notes"`
}

type ItemReq struct {
	VaultId  int
	ItemName string
	Username string
	Password string
	Url      string
	Notes    string
}

func main() {
	db.InitDB("postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	defer db.CloseDB()

	http.Handle("/save-password", authMiddleware(http.HandlerFunc(savePassword)))
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func savePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(ItemReq)
	err := body.Decode(&params)

	vaultid := params.VaultId
	itemName := params.ItemName
	username := params.Username
	password := params.Password
	url := params.Url
	notes := params.Notes

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	// Create a new item
	newItem := db.Item{VaultId: vaultid, ItemName: itemName, UserName: username, EncryptedPassword: hash, Url: url, Notes: notes}
	dberr := db.CreateItem(newItem)
	if dberr != nil {
		fmt.Println("Error creating item:", err)
	}

	json.NewEncoder(w).Encode(params)
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

func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username, password := auth.DecodeBasicAuthHeader(authHeader)
		if !auth.BasicAuth(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

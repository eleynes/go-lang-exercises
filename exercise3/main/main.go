package main

import (
	"encoding/hex"
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

type UserReq struct {
	UserName string
	Email    string
	Password string
	Salt     string
}

type ValultReq struct {
	VaultName string
}
type ItemReq struct {
	VaultName string
	ItemName  string
	Username  string
	Password  string
	Url       string
	Notes     string
}
type GenerateItemReq struct {
	VaultName           string
	ItemName            string
	Username            string
	Url                 string
	Notes               string
	Length              int
	NumPasswords        int
	PasswordType        string
	IsNumbersIncluded   bool
	IsSymbolsIncluded   bool
	IsUppercaseIncluded bool
}

func main() {
	db.InitDB("postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	defer db.CloseDB()

	// public
	// http.HandleFunc("/create-user", createUser)
	http.HandleFunc("/generate-password", generatePassword)

	// private
	// http.Handle("/create-user", authMiddleware(http.HandlerFunc(createUser)))
	http.Handle("/create-vault", auth.AuthMiddleware(http.HandlerFunc(createVault)))
	http.Handle("/create-item", auth.AuthMiddleware(http.HandlerFunc(createItem)))
	http.Handle("/generate-items", auth.AuthMiddleware(http.HandlerFunc(generateItem)))
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(UserReq)
	err := body.Decode(&params)

	userName := params.UserName
	email := params.Email
	password := params.Password
	salt := params.Salt

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	hash := string(bytes)

	// Create a new item
	newUser := db.User{UserName: userName, Email: email, MasterpasswordHash: hash, MasterpasswordSalt: salt}
	dberr := db.CreateUser(newUser)
	if dberr != nil {
		fmt.Println("Error creating item:", err)
	}

	json.NewEncoder(w).Encode(newUser)
}

func createVault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(ValultReq)
	err := body.Decode(&params)

	vaultName := params.VaultName

	// Create a new vault
	newVault := db.Vault{VaultName: vaultName, UserId: auth.CurrentLoggedInUser.UserId}
	dberr := db.CreateVault(newVault)
	if dberr != nil {
		fmt.Println("Error creating item:", err)
	}
	vault, err := db.GetVaultByNameAndUserId(vaultName, auth.CurrentLoggedInUser.UserId)

	json.NewEncoder(w).Encode(vault)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(ItemReq)
	err := body.Decode(&params)

	vaultName := params.VaultName
	itemName := params.ItemName
	username := params.Username
	password := params.Password
	url := params.Url
	notes := params.Notes

	key := []byte("16byteAESKey1234")
	message := []byte(password)

	encrypted := db.Encrypt(message, key)
	fmt.Println("Encrypted:", hex.EncodeToString(encrypted))

	// decrypted := db.Decrypt(encrypted, key)
	// fmt.Println("Decrypted:", string(decrypted))

	// Create a new item
	// TODO: validation for vault ownership
	newItem := db.Item{VaultName: vaultName, UserId: auth.CurrentLoggedInUser.UserId, ItemName: itemName, UserName: username, EncryptedPassword: hex.EncodeToString(encrypted), Url: url, Notes: notes}
	itemid, dberr := db.CreateItem(newItem)
	if dberr != nil {
		fmt.Println("Error creating item:", err)
	}

	item, err := db.GetItemByItemId(itemid)
	item.UserId = auth.CurrentLoggedInUser.UserId
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(item)
}

func generateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(GenerateItemReq)
	err := body.Decode(&params)

	var passwords []Password
	var items []db.Item
	passwords, err = generatePasswords(params.NumPasswords, params.PasswordType, params.Length, params.IsNumbersIncluded, params.IsSymbolsIncluded, params.IsUppercaseIncluded)

	for i := 0; i < len(passwords); i++ {
		key := []byte("16byteAESKey1234")
		message := []byte(passwords[i].Password)

		encrypted := db.Encrypt(message, key)
		fmt.Println("Encrypted:", hex.EncodeToString(encrypted))

		// decrypted := db.Decrypt(encrypted, key)
		// fmt.Println("Decrypted:", string(decrypted))

		// Create a new item
		newItem := db.Item{VaultName: params.VaultName, UserId: auth.CurrentLoggedInUser.UserId, ItemName: params.ItemName, UserName: params.Username, EncryptedPassword: hex.EncodeToString(encrypted), Url: params.Url, Notes: params.Notes}
		itemid, dberr := db.CreateItem(newItem)
		if dberr != nil {
			fmt.Println("Error creating item:", err)
		}

		item, err := db.GetItemByItemId(itemid)
		item.UserId = auth.CurrentLoggedInUser.UserId

		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func generatePassword(w http.ResponseWriter, r *http.Request) {
	var passwords []Password
	w.Header().Set("Content-Type", "application/json")

	body := json.NewDecoder(r.Body)
	params := new(PasswordReq)
	err := body.Decode(&params)

	passwords, err = generatePasswords(params.NumPasswords, params.PasswordType, params.Length, params.IsNumbersIncluded, params.IsSymbolsIncluded, params.IsUppercaseIncluded)

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(passwords)
}

func generatePasswords(numPasswordsInt int, passwordType string, lengthInt int, includeNumbersBool bool, includeSymbolsBool bool, includeUppercaseBool bool) ([]Password, error) {
	var passwords []Password
	var err error
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
			return passwords, err
		}

		if err != nil {
			fmt.Println("Error generating password/PIN:", err)
			return passwords, err
		}

		var passwordobj Password
		passwordobj.ID = i + 1
		passwordobj.Password = password

		passwords = append(passwords, passwordobj)
	}
	return passwords, nil
}

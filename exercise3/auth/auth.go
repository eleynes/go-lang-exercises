package auth

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"example.com/db"
	"golang.org/x/crypto/bcrypt"
)

func BasicAuth(username, password string) bool {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	// decodedPwd, err := base64.StdEncoding.DecodeString(password)
	// // decodedPwd, err := base64.StdEncoding.E(password)
	fmt.Println("HASH:", hash)
	// fmt.Println("User:", string(decodedPwd))
	user, err := db.GetUserByUsername(username)
	if err != nil {
		log.Fatal(err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.MasterpasswordHash), []byte(password))
	if err != nil {
		fmt.Println("Error comparing passwords:", err)
		return false
	}

	return true

}

func DecodeBasicAuthHeader(authHeader string) (string, string) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", ""
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ""
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", ""
	}

	return creds[0], creds[1]
}

package auth

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"example.com/db"
	"golang.org/x/crypto/bcrypt"
)

type CurrentUser struct {
	UserId             int
	UserName           string
	Email              string
	MasterpasswordHash string
	MasterpasswordSalt string
	CreatedAt          string
	UpdatedAt          string
}

var CurrentLoggedInUser CurrentUser

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username, password := decodeBasicAuthHeader(authHeader)
		if !basicAuth(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func basicAuth(username, password string) bool {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		log.Fatal(err)
		return false
	}

	CurrentLoggedInUser = CurrentUser(user)

	err = bcrypt.CompareHashAndPassword([]byte(user.MasterpasswordHash), []byte(password+user.MasterpasswordSalt))
	if err != nil {
		fmt.Println("Error comparing passwords:", err)
		return false
	}

	return true

}

func decodeBasicAuthHeader(authHeader string) (string, string) {
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

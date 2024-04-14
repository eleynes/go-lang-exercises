package db

import (
	"log"
)

type User struct {
	UserId             int
	UserName           string
	Email              string
	MasterpasswordHash string
	MasterpasswordSalt string
	CreatedAt          string
	UpdatedAt          string
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT userid, username, email, masterpasswordhash, masterpasswordsalt, createdat, updatedat FROM users WHERE username = $1", username).Scan(&user.UserId, &user.UserName, &user.Email, &user.MasterpasswordHash, &user.MasterpasswordSalt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}
	return user, nil
}

package db

import (
	"log"
)

type Item struct {
	ItemId            int
	VaultId           int
	ItemName          string
	UserName          string
	EncryptedPassword string
	Url               string
	Notes             string
	CreatedAt         string
	UpdatedAt         string
}

func CreateItem(item Item) error {
	_, err := db.Exec("INSERT INTO items (vaultid, itemname, username, encryptedpassword, url, notes) VALUES ($1, $2, $3, $4, $5, $6)", item.VaultId, item.ItemName, item.UserName, item.EncryptedPassword, item.Url, item.Notes)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

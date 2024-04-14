package db

import (
	"log"
)

type Vault struct {
	vaultid   int
	vaultname string
	userid    int
	createdat string
	updatedat string
}

func CreateVault(vault Vault) error {
	_, err := db.Exec("INSERT INTO vaults (vaultname, userid) VALUES ($1, $2)", vault.vaultname, vault.userid)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

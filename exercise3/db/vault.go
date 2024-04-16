package db

import (
	"log"
)

type Vault struct {
	VaultId   int
	VaultName string
	UserId    int
	CreatedAt string
	UpdatedAt string
}

func CreateVault(vault Vault) error {
	_, err := db.Exec("INSERT INTO vaults (vaultname, userid) VALUES ($1, $2)", vault.VaultName, vault.UserId)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetVaultByNameAndUserId(name string, userid int) (Vault, error) {
	var vault Vault
	err := db.QueryRow("SELECT vaultid, vaultname, userid, createdat, updatedat FROM vaults WHERE vaultname = $1 AND userid = $2", name, userid).Scan(&vault.VaultId, &vault.VaultName, &vault.UserId, &vault.CreatedAt, &vault.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return Vault{}, err
	}
	return vault, nil
}

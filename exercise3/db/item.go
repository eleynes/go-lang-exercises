package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

type Item struct {
	ItemId            int
	VaultId           int
	VaultName         string
	ItemName          string
	UserName          string
	EncryptedPassword string
	Url               string
	Notes             string
	CreatedAt         string
	UpdatedAt         string
}

func CreateItem(item Item) error {
	_, err := db.Exec("INSERT INTO items (vaultid, itemname, username, encryptedpassword, url, notes) VALUES ((SELECT vaultid FROM vaults WHERE vaultname = $1), $2, $3, $4, $5, $6)", item.VaultName, item.ItemName, item.UserName, item.EncryptedPassword, item.Url, item.Notes)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetItemByNameAndVaultId(name string, vaultid int) (Item, error) {
	var item Item
	err := db.QueryRow("SELECT itemid, vaults.vaultid, vaultname, itemname, username, encryptedpassword, url, notes, items.createdat, items.updatedat FROM items LEFT JOIN vaults ON items.vaultid = vaults.vaultid WHERE itemname = $1 AND items.vaultid = $2", name, vaultid).Scan(&item.ItemId, &item.VaultId, &item.VaultName, &item.ItemName, &item.UserName, &item.EncryptedPassword, &item.Url, &item.Notes, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return Item{}, err
	}
	return item, nil
}

func Encrypt(plaintext []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

func Decrypt(ciphertext []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext
}

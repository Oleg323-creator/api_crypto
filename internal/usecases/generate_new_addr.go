package usecases

import (
	"api_crypto/internal/db/repository"
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func (u *Usecases) GenerateNewAdd(curr string) (string, error) {
	err := godotenv.Load("/home/oleh/GolandProjects/api_crypto/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicAddres := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	fmt.Printf("Public address: %s", publicAddres)

	// CONVERT KEY TO DER FORMAT
	derKey := crypto.FromECDSA(privateKey)

	password := []byte(os.Getenv("SECRET_PASSWORD"))
	if os.Getenv("SECRET_PASSWORD") == "" {
		log.Fatal("Error getting secret password from .env")
	}
	log.Printf("Using secret password: %s", os.Getenv("SECRET_PASSWORD"))

	encryptedKey, err := u.EncryptAESGCM(password, derKey)
	if err != nil {
		log.Fatal(err)
	}

	encryptedKeyBase64 := base64.StdEncoding.EncodeToString(encryptedKey)

	data := repository.DataToSave{
		PrivateKey: encryptedKeyBase64,
		Address:    publicAddres,
		Currency:   curr,
	}

	err = u.Repository.SaveNewAddrToDB(data)
	if err != nil {
		return "", fmt.Errorf("Error saiving new addr to DB: %v", err)
	}

	return publicAddres, nil
}

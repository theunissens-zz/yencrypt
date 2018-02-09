package encryptserver

import (
	"errors"
	. "github.com/yencrypt/yencrypt_server/encryptserver/crypto"
	. "github.com/yencrypt/yencrypt_server/encryptserver/db"
	"github.com/yencrypt/yencrypt_server/encryptserver/exceptions"
	"log"
)

type YEncryptServiceInterface interface {
	ProcessEncryption(id string, plainText string) (string, error)
	ProcessDecryption(id string, key string) (string, error)
}

type YEncryptService struct {
	DB      YDatabaseInterface
	YCrypto YCrypto
}

func (service *YEncryptService) init() error {
	service.DB = &YDatabase{
		DBName:   DBName,
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
	}
	err := service.DB.Setup()
	if err != nil {
		log.Printf("An error occurred setting up service: %v", err.Error())
		return err
	}
	err = service.DB.Connect(false)
	if err != nil {
		log.Printf("An error occurred setting up service: %v", err.Error())
		return err
	}
	return nil
}

func (service *YEncryptService) ProcessEncryption(id string, plainText string) (string, error) {
	// Convert to byte array
	if len(plainText) == 0 {
		log.Print("Plaintext cannot be empty")
		return "", &exceptions.ValidationError{ErrorString: "Plaintext cannot be empty"}
	}

	key, err := service.YCrypto.GenerateKey()
	if err != nil {
		log.Printf("Error generating the key: %v", err.Error())
		return "", err
	}
	cypherText, err := service.YCrypto.Encrypt(key, []byte(plainText))
	if err != nil {
		log.Printf("An error occurred trying to encrypt data: %v", err.Error())
		return "", err
	}

	encodedCypherText := service.YCrypto.ConvertToBase64(cypherText)

	hashId := service.YCrypto.Hash(id)

	if err := service.DB.Store(hashId, encodedCypherText); err != nil {
		log.Printf("Unable to store encrypted data: %v", err.Error())
		return "", err
	}

	return service.YCrypto.ConvertToBase64(key), nil
}

func (service *YEncryptService) ProcessDecryption(id string, b64Key string) (string, error) {
	if len(b64Key) == 0 {
		return "", errors.New("Key cannot be empty")
	}

	var err error
	key, err := service.YCrypto.ConvertFromBase64(b64Key)
	if err != nil {
		log.Printf("Error decoding key: %v", id)
		return "", err
	}

	hashId := service.YCrypto.Hash(id)

	cypherText, err := service.retrieveCypherText(hashId)
	if err != nil {
		log.Printf("Could not retrieve cypher text for id: %v", id)
		return "", err
	}

	cypherTextData, err := service.YCrypto.ConvertFromBase64(cypherText)
	if err != nil {
		log.Printf("Error decoding cyphertext: %v", id)
		return "", err
	}

	plainText, err := service.YCrypto.Decrypt(key, cypherTextData)
	if err != nil {
		log.Printf("Error decrypting cypherdata: %v", id)
		return "", err
	}

	return string(plainText[0:]), nil
}

func (service *YEncryptService) retrieveCypherText(id string) (string, error) {
	cypherYext, err := service.DB.Retrieve(id)
	if err != nil {
		return "", err
	}
	return cypherYext, nil
}

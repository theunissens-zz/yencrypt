package test

import (
	"bytes"
	. "github.com/yencrypt/yencrypt_server/encryptserver/crypto"
	"strings"
	"testing"
)

func TestKeyGeneration_success(t *testing.T) {
	yCrypto := YCrypto{}
	key, err := yCrypto.GenerateKey()
	if err != nil {
		t.Error("Error generating key")
	}
	if len(key) != KeyLength {
		t.Error("Key generate not expected length")
	}
}

func TestEncryption_success(t *testing.T) {
	plainText := "someRandomText"
	plainTextArr := []byte(plainText)

	// Encrypt plaintext
	yCrypto := YCrypto{}
	key, err := yCrypto.GenerateKey()
	if err != nil {
		t.Error("Error generating key")
	}
	cypherText, err := yCrypto.Encrypt(key, plainTextArr)
	if err != nil {
		t.Errorf("An unexpected error occurred trying to encrypt plaintext: %s error: %v", plainText, err)
	}
	if cypherText == nil || bytes.Equal(cypherText, plainTextArr) {
		t.Errorf("An unexpected error occurred trying to encrypt plaintext: %s error: %v", plainText, err)
	}
}

func TestDecryption_success(t *testing.T) {
	strPlainText := "someRandomText"
	plainText := []byte(strPlainText)

	// Encrypt plaintext
	yCrypto := YCrypto{}
	key, err := yCrypto.GenerateKey()
	if err != nil {
		t.Error("Error generating key")
	}

	cypherText, err := yCrypto.Encrypt(key, plainText)
	if err != nil {
		t.Errorf("An unexpected error occurred trying to encrypt plaintext: %s error: %v", strPlainText, err)
	}
	if cypherText == nil || bytes.Equal(cypherText, plainText) {
		t.Errorf("An unexpected error occurred trying to encrypt plaintext: %s error: %v", strPlainText, err)
	}

	// Convert to and from base64
	c64CypherText := yCrypto.ConvertToBase64(cypherText)
	c64CypherTextAfter, _ := yCrypto.ConvertFromBase64(string(c64CypherText[0:]))

	// Decrypt cyphertext
	decryptedPlainText, err := yCrypto.Decrypt(key, []byte(c64CypherTextAfter))
	if err != nil {
		t.Errorf("An error occurred trying to decrpyt, error: %v", err)
	}

	strDecryptedPlainText := string(decryptedPlainText[0:])
	if strings.Compare(strDecryptedPlainText, strPlainText) != 0 {
		t.Error("Decrypted plaintext is not the same as original plaintext")
	}
}

func TestHashId_success(t *testing.T) {
	id := "someId"
	yCrypto := YCrypto{}
	idHash := yCrypto.Hash(id)
	if len(idHash) == 0 {
		t.Error("Hash should not be == 0")
	}
}

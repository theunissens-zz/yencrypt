package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

const KeyLength = 32

type YCrypto struct {
}

func (c *YCrypto) GenerateKey() ([]byte, error) {
	key := make([]byte, KeyLength)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (c *YCrypto) Encrypt(key, plainText []byte) (cypherText []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	cypherText = make([]byte, aes.BlockSize+len(plainText))
	iv := cypherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cypherText[aes.BlockSize:], plainText)
	return
}

func (c *YCrypto) Decrypt(key, cypherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cypherText) < aes.BlockSize {
		return nil, errors.New("Length of cypher text needs to be the same as block size")
	}

	iv := cypherText[:aes.BlockSize]
	cypherText = cypherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cypherText, cypherText)

	return cypherText, nil
}

func (c *YCrypto) ConvertToBase64(data []byte) string {
	b64Data := base64.StdEncoding.EncodeToString(data)
	return b64Data
}

func (c *YCrypto) ConvertFromBase64(data string) ([]byte, error) {
	res, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Hashes and converts to base64
func (c *YCrypto) Hash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	hashData := hash.Sum(nil)
	return c.ConvertToBase64(hashData)
}

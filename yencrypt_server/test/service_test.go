package test

import (
	"encoding/base64"
	. "github.com/yencrypt/yencrypt_server/encryptserver"
	"strings"
	"testing"
)

func TestServiceEncrypt_success(t *testing.T) {
	service := YEncryptService{}
	service.DB = &MockYDatabase{}
	_, err := service.ProcessEncryption("1", "somePlainText")
	if err != nil {
		t.Errorf("An unexpected error occurred: %v: ", err)
	}
}

func TestServiceEncrypt_noPlainText_fail(t *testing.T) {
	service := YEncryptService{}
	service.DB = &MockYDatabase{}
	_, err := service.ProcessEncryption("1", "")
	if !strings.Contains(err.Error(), "Plaintext cannot be empty") {
		t.Error("Expected error as no plaintext was supplied")
	}
}

func TestServiceDecrypt_noKey_fail(t *testing.T) {
	service := YEncryptService{}
	service.DB = &MockYDatabase{}
	var key string
	_, err := service.ProcessDecryption("1", key)
	if !strings.Contains(err.Error(), "Key cannot be empty") {
		t.Error("Expected error as no key was supplied")
	}
}

type MockYDatabase struct {
}

func (db *MockYDatabase) Setup() error {
	return nil
}

func (db *MockYDatabase) Connect(connectToServer bool) error {
	return nil
}

func (service *MockYDatabase) Store(id string, data string) error {
	return nil
}

func (service *MockYDatabase) Retrieve(id string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte("someCypherText")), nil
}

func (db *MockYDatabase) DropDB() error {
	return nil
}

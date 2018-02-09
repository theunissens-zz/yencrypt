package test

import (
	"bytes"
	"fmt"
	. "github.com/yencrypt/yencrypt_server/encryptserver"
	"github.com/yencrypt/yencrypt_server/encryptserver/exceptions"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncryptionPost_Success(t *testing.T) {
	yServer := YServer{}
	service := MockYEncryptionServiceSuccess{}
	yServer.Service = &service

	handler := &EncryptDataHandler{Server: &yServer}
	server := httptest.NewServer(handler)
	defer server.Close()

	body := "justSomeRandom\nText"

	resp, err := http.Post(server.URL+"/encrypt/someId", "text/plain", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Error("Could not make post to encryption server")
	}
	if resp.StatusCode != 200 {
		t.Errorf("Received %v from server, expected 200", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body %v", err)
	}
	strResponseBody := string(responseBody)
	if len(strResponseBody) == 0 {
		t.Error("No body received")
	}
}

func TestEncryptionPost_WrongHttpMethod_Fail(t *testing.T) {
	server := httptest.NewServer(&EncryptDataHandler{})
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error("Could not make post to encryption server")
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Received %v from server, expected %v", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestEncryptionPost_ValidationError_Fail(t *testing.T) {
	yServer := YServer{}
	service := MockYEncryptionServiceFailure{}
	yServer.Service = &service

	handler := &EncryptDataHandler{Server: &yServer}
	server := httptest.NewServer(handler)
	defer server.Close()

	body := "somePlainText"

	resp, err := http.Post(server.URL+"/encrypt/someId", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Error("Could not make post to encryption server")
	}
	if resp.StatusCode != 400 {
		t.Errorf("Received %v from server, expected 400", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body %v", err)
	}
	strResponseBody := string(responseBody)
	fmt.Print(strResponseBody)
	if !strings.Contains(strResponseBody, "Some validation error") {
		t.Error("Incorrect error received")
	}
}

func TestDecryptionGet_Success(t *testing.T) {
	yServer := YServer{}
	service := MockYEncryptionServiceSuccess{}
	yServer.Service = &service

	handler := &DecryptDataHandler{Server: &yServer}
	server := httptest.NewServer(handler)
	defer server.Close()

	path := fmt.Sprintf("/%s/%d?key=%s", "encrypt", 1, "someKey")

	resp, err := http.Get(server.URL + path)
	if err != nil {
		t.Errorf("Could not make post to encryption server: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Received %v from server, expected 200", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body %v", err)
	}
	strResponseBody := string(responseBody)
	if len(strResponseBody) == 0 {
		t.Error("No body received")
	}
	if strings.Compare(strResponseBody, "somePlainText") != 0 {
		t.Error("Incorrect plaintext received")
	}
}

type MockYEncryptionServiceSuccess struct {
}

func (service *MockYEncryptionServiceSuccess) ProcessEncryption(id string, plainText string) (string, error) {
	return "someKey", nil
}

func (service *MockYEncryptionServiceSuccess) ProcessDecryption(id string, key string) (string, error) {
	return "somePlainText", nil
}

type MockYEncryptionServiceFailure struct {
}

func (service *MockYEncryptionServiceFailure) ProcessEncryption(id string, plainText string) (string, error) {
	return "", &exceptions.ValidationError{"Some validation error"}
}

func (service *MockYEncryptionServiceFailure) ProcessDecryption(id string, key string) (string, error) {
	return "", &exceptions.ValidationError{"Some validation error"}
}

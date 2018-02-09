package test

import (
	"bytes"
	"fmt"
	. "github.com/yencrypt/yencrypt_server/encryptserver"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
)

const testPort = 54321

func TestEndToEnd_success(t *testing.T) {
	serverSyncChannel := make(chan bool)
	yServer := &YServer{}

	go startUp(yServer, t)
	blockTillServerIsUp()
	clients := func() {
		var wg sync.WaitGroup
		wg.Add(100)
		for i := 0; i < 100; i++ {
			go sendRequests(&wg, i, t)
		}
		wg.Wait()
		serverSyncChannel <- true
	}
	go clients()

	stop := <-serverSyncChannel

	if stop {
		// We want a fresh db every time
		service, _ := yServer.Service.(*YEncryptService)
		err := service.DB.DropDB()
		if err != nil {
			t.Error("Error occurred dropping db: ", err)
		}
	}
}

func blockTillServerIsUp() {
	started := false
	for !started {
		_, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", testPort))
		if err == nil {
			started = true
		}
	}
}

func sendRequests(wg *sync.WaitGroup, clientNumber int, t *testing.T) {
	for i := 0; i < 100; i++ {
		someId := fmt.Sprintf("someId %d %d", clientNumber, i)
		somePlainText := fmt.Sprintf("Here is some cool plaintext to encrypt %d %d", clientNumber, i)

		key := sendEncryptReq(someId, somePlainText, t)

		plainTextResult := sendDecryptReq(someId, key, t)

		if strings.Compare(somePlainText, plainTextResult) != 0 {
			t.Errorf("Received plaintext incorrect: %s", plainTextResult)
		}
	}
	wg.Done()
}

func sendEncryptReq(id, plainText string, t *testing.T) string {
	strReq := fmt.Sprintf("http://localhost:%d/encrypt/%s", testPort, id)
	resp, err := http.Post(strReq, "text/plain", bytes.NewBuffer([]byte(plainText)))
	if err != nil {
		t.Fatalf("Error sending encrypt request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got: %d", resp.StatusCode)
	}
	key, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("An error occurred reading response body: %v", err)
	}
	return string(key[0:])
}

func sendDecryptReq(id, key string, t *testing.T) string {
	url := fmt.Sprintf("http://localhost:%d/decrypt/%s?key=%s", testPort, id, url.QueryEscape(string(key[0:])))
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Error sending encrypt request: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got: %d", resp.StatusCode)
	}
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("An error occurred reading response body: %v", err)
	}
	return string(payload[0:])
}

func startUp(yServer *YServer, t *testing.T) error {
	err := yServer.StartServer(testPort)
	if err != nil {
		t.Fatal("Server failed to start up")
	}
	return nil
}

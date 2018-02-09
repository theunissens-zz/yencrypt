package client

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var Port = 12345

// Client provides functionality to interact with the encryption-server
type Client interface {
	// Store accepts an id and a payload in bytes and requests that the
	// encryption-server stores them in its data store
	Store(id, payload []byte) (aesKey []byte, err error)

	// Retrieve accepts an id and an AES key, and requests that the
	// encryption-server retrieves the original (decrypted) bytes stored
	// with the provided id
	Retrieve(id, aesKey []byte) (payload []byte, err error)
}

type YClient struct {
}

func (c *YClient) Store(id, payload []byte) (aesKey []byte, err error) {
	strReq := fmt.Sprintf("http://localhost:%d/encrypt/%s", Port, url.QueryEscape(string(id[0:])))
	resp, err := http.Post(strReq, "text/plain", bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		errorResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("An error occurred reading response body: %v", err)
			return nil, err
		} else {
			err = errors.New(fmt.Sprintf("An error occurred posting encryption data: %v", string(errorResponse[0:])))
			return nil, err
		}
	}
	aesKey, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An error occurred reading response body: %v", err)
	}
	return
}

func (c *YClient) Retrieve(id, aesKey []byte) (payload []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/decrypt/%s?key=%s", Port, url.QueryEscape(string(id[0:])), url.QueryEscape(string(aesKey[0:]))))
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		errorResponse, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("An error occurred reading response body: %v", err)
			return nil, err
		} else {
			err = errors.New(fmt.Sprintf("An error occurred requesting decryption: %v", string(errorResponse[0:])))
			return nil, err
		}
	}
	payload, err = ioutil.ReadAll(resp.Body)
	return
}

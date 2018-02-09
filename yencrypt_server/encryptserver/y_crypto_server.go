package encryptserver

import (
	"fmt"
	"github.com/yencrypt/yencrypt_server/encryptserver/exceptions"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type YServer struct {
	Service YEncryptServiceInterface
	Server  http.Server
}

func (yServer *YServer) StartServer(port int) error {
	service := YEncryptService{}
	err := service.init()
	if err != nil {
		return err
	}
	yServer.Service = &service

	// Create handlers
	pingHandler := &PingHandler{}
	encryptDataHandler := &EncryptDataHandler{Server: yServer}
	decryptDataHandler := &DecryptDataHandler{Server: yServer}

	mux := http.NewServeMux()
	mux.Handle("/ping", pingHandler)
	mux.Handle("/encrypt/", encryptDataHandler)
	mux.Handle("/decrypt/", decryptDataHandler)
	s := http.Server{Handler: mux, Addr: fmt.Sprintf("localhost:%d", port)}
	yServer.Server = s
	err = s.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

type PingHandler struct {
}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Incorrect http method used for request url %s", r.URL)
		http.Error(w, "Incorrect http method used", http.StatusMethodNotAllowed)
		return
	}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Fprintf(w, "Service is alive: %v", timestamp)

}

type EncryptDataHandler struct {
	Server *YServer
}

func (h *EncryptDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validation
	if r.Method != http.MethodPost {
		log.Printf("Incorrect http method used for request url %s", r.URL)
		http.Error(w, "Incorrect http method used", http.StatusMethodNotAllowed)
		return
	}

	// Parse path params
	path := r.URL.Path
	values := strings.Split(path, "/")

	// Get id
	strId := values[2]

	// Get plaintext
	plainText, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Unable to read body from request", http.StatusBadRequest)
		return
	}

	// Pass to service layer for processing
	key, err := h.Server.Service.ProcessEncryption(strId, string(plainText[0:]))
	if err != nil {
		log.Printf("Error processing encryption request: %v", err)
		_, ok := err.(*exceptions.ValidationError)
		if ok {
			http.Error(w, fmt.Sprintf("Error processing encryption request: %v", err.Error()), http.StatusBadRequest)
			return
		} else {
			http.Error(w, fmt.Sprintf("Error processing encryption request: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		return
	}
	fmt.Fprint(w, key)
}

type DecryptDataHandler struct {
	Server *YServer
}

func (h *DecryptDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Incorrect http method used for request url %s", r.URL)
		http.Error(w, "Incorrect http method used", http.StatusMethodNotAllowed)
		return
	}

	// Parse path params
	path := r.URL.Path
	path = strings.TrimPrefix(path, "/decrypt/")
	values := strings.SplitN(path, "/", 1)
	// Get id
	strId := values[0]

	// Get key from query param
	strKey := r.URL.Query().Get("key")
	if len(strKey) == 0 {
		http.Error(w, "Key cannot be empty", http.StatusBadRequest)
		return
	}

	// Pass to service layer for processing
	plainText, err := h.Server.Service.ProcessDecryption(strId, strKey)
	if err != nil {
		log.Printf("Error processing decryption request: %v", err)
		_, ok := err.(*exceptions.ValidationError)
		if ok {
			http.Error(w, fmt.Sprintf("Error processing decryption request: %v", err.Error()), http.StatusBadRequest)
			return
		} else {
			http.Error(w, fmt.Sprintf("Error processing decryption request: %v", err.Error()), http.StatusInternalServerError)
			return
		}
		return
	}
	fmt.Fprint(w, plainText)
}

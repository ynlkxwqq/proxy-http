// package main

// import "simple-http-proxy/internal/app"

// func main() {
// 	app.Run()
// }

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Структуры для входящего запроса и ответа
type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type ProxyResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int               `json:"length"`
}

// Хранилище запросов и ответов
var storage sync.Map

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Method == "" || req.URL == "" {
		http.Error(w, "Missing method or URL", http.StatusBadRequest)
		return
	}

	// Формируем новый HTTP-запрос к стороннему сервису
	clientReq, err := http.NewRequest(req.Method, req.URL, nil)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for k, v := range req.Headers {
		clientReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(clientReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error contacting target: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		respHeaders[k] = v[0]
	}

	id := uuid.New().String()
	proxyResp := ProxyResponse{
		ID:      id,
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Length:  len(body),
	}

	// Сохраняем запрос и ответ
	storage.Store(id, map[string]interface{}{
		"request":  req,
		"response": proxyResp,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxyResp)
}

func main() {
	http.HandleFunc("/proxy", proxyHandler)
	fmt.Println("Server listening on :3333")
	log.Fatal(http.ListenAndServe(":3333", nil))
}

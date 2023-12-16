package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RESTServer struct {
	server    *http.Server
	OllamaURL string
}

type OllamaRequestGenerate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaRequestPull struct {
	Name   string `json:"name"`
	Stream bool   `json:"stream"`
}

func NewRESTServer(address, ollamaURL string) *RESTServer {
	return &RESTServer{
		server: &http.Server{
			Addr:    address,
			Handler: mux.NewRouter(),
		},
		OllamaURL: ollamaURL,
	}
}

func (s *RESTServer) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/generate", s.generate).Methods("POST")
	r.HandleFunc("/api/pull", s.pull).Methods("POST")

	s.server.Handler = r

	log.Printf("Starting REST server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *RESTServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Stopping REST server gracefully...")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	} else {
		log.Println("REST server stopped")
	}
}

func (s *RESTServer) generate(w http.ResponseWriter, r *http.Request) {
	var req OllamaRequestGenerate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := s.forwardToOllama(req, "/api/generate")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, response)
}

func (s *RESTServer) pull(w http.ResponseWriter, r *http.Request) {
	var req OllamaRequestPull
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := s.forwardToOllama(req, "/api/pull")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, response)
}

func (s *RESTServer) forwardToOllama(req interface{}, url string) (string, error) {
	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(s.OllamaURL+url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

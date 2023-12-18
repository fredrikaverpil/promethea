package servers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type OllamaRequestGenerate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaRequestPull struct {
	Name   string `json:"name"`
	Stream bool   `json:"stream"`
}

type OllamaResponseGenerate struct {
	Response string `json:"response"`
}

type RequestErrors struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Value   string `json:"value"`
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

func (s *RESTServer) errorsGuessField(w http.ResponseWriter, r *http.Request) {
	var req RequestErrors
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read YAML of instructions
	instructionsFile := "instructions/errors.yaml"
	instructionsPromptBytes, err := ioutil.ReadFile(instructionsFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	prompt := string(instructionsPromptBytes)

	// Extend instructionsPrompt with RequestErrors
	prompt = prompt + "\n\n"
	prompt = prompt + "Original error message: `" + req.Message + "`\n"
	prompt = prompt + "Original error code: `" + req.Code + "`\n"
	prompt = prompt + "Original error value: `" + req.Value + "`\n"

	// print out prompt with linebreaks
	fmt.Println(prompt)

	reqOllama := OllamaRequestGenerate{
		Model:  "mistral",
		Prompt: prompt,
		Stream: false,
	}

	respOllama, err := s.forwardToOllama(reqOllama, "/api/generate")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprint(w, response)

	// marshal response into OllamaResponseGenerate struct
	var resp OllamaResponseGenerate
	if err := json.Unmarshal([]byte(respOllama), &resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set header
	w.Header().Set("Content-Type", "application/json")

	// write response to response writer
	fmt.Fprint(w, resp.Response)
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

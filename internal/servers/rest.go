package servers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RESTServer struct {
	server    *http.Server
	OllamaURL string
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
	r.HandleFunc("/api/errors/guess_field", s.errorsGuessField).Methods("POST")

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

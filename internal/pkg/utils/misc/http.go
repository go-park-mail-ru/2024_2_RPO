package misc

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func StartServers(servers ...*http.Server) {
	quit := make(chan os.Signal, 1)
	errChan := make(chan error, len(servers))
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запустить серверы
	for _, server := range servers {
		go func(s *http.Server) {
			log.Printf("Starting server on %s\n", s.Addr)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}(server)
	}

	select {
	case sig := <-quit:
		log.Printf("Received signal %s. Shutting down...\n", sig)
		os.Exit(0)
	case err := <-errChan:
		log.Fatalf("Server error: %v\n", err)
	}
}

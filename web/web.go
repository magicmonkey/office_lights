package web

import (
	"embed"
	"log"
	"net/http"
	"sync"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the web server
type Server struct {
	ledStrip    *ledstrip.LEDStrip
	ledBar      *ledbar.LEDBar
	videoLight1 *videolight.VideoLight
	videoLight2 *videolight.VideoLight
	httpServer  *http.Server
	mu          sync.Mutex // Protect concurrent access
}

// NewServer creates a new web server
func NewServer(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) *Server {
	return &Server{
		ledStrip:    strip,
		ledBar:      bar,
		videoLight1: vl1,
		videoLight2: vl2,
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	// Serve index.html at root
	mux.HandleFunc("/", s.handleIndex)

	// API endpoints
	mux.HandleFunc("/api", s.handleAPI)
	mux.HandleFunc("/health", s.handleHealth)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Web server starting on http://localhost:%s", port)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}

// handleIndex serves the main HTML page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Only serve index.html for the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

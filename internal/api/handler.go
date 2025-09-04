package api

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	rtp    float64
}

func NewServer(logger *zap.Logger, rtp float64) *Server {
	return &Server{
		logger: logger,
		rtp:    rtp,
	}
}

func (s *Server) RegisterRoutes() *http.ServeMux {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /get", s.HandleGet)
	return serveMux
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	multiplier := 1.0 + rand.Float64()*(10000.0-1.0) // TODO

	resp := map[string]string{
		"result": strconv.FormatFloat(multiplier, 'f', 4, 64),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

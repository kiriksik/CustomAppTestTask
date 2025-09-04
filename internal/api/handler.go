package api

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	logger               *zap.Logger
	rtp                  float64
	maxGenerate          float64
	maxClientGenerate    float64
	count                float64
	sumExpectedGenerated float64
}

func NewServer(logger *zap.Logger, rtp float64) *Server {
	return &Server{
		logger:               logger,
		rtp:                  rtp,
		maxGenerate:          500,
		maxClientGenerate:    10000,
		count:                0,
		sumExpectedGenerated: 0,
	}
}

func (s *Server) RegisterRoutes() *http.ServeMux {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /get", s.HandleGet)
	return serveMux
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	// если N запросов для чисел от 0 до 10 000,
	// то среднее число в каждом запросе 5 000
	// sum0 = N
	// sum1 = 5 000 * N
	// значит, если RTP допустим 0.95, то sum1 должно быть
	// sum1expected = 0.95 * 5 000 * И = 4 750 * N
	// т.е. фактически, от генератора приходит 500 * N / 10 000 чисел до 500
	// т.е. в среднем 250 * N / 10 000 = 125 000 для чисел до 500
	// значит при ограничении генератора 500,
	// при увеличении количества запросов точность (приближение RES к RTP) будет расти
	// значит в среднем 50% чисел до 500 выживут,
	// их среднее значение 250,
	// т.е. в среднем sum1 = 62 500

	var multiplier float64
	var p_zero float64
	var expectedRES float64
	if s.count == 0 {
		expectedRES = 0
	} else {
		expectedRES = s.sumExpectedGenerated / s.count
	}
	if expectedRES > s.rtp {
		p_zero = 1
	} else {
		meanGenerator := s.maxGenerate / 2
		p_zero = 1 - (((s.count+1)*s.rtp - s.sumExpectedGenerated) / meanGenerator)
	}

	multiplier = rand.Float64() * (s.maxGenerate)
	if rand.Float64() < p_zero {
		multiplier = 0
	}
	s.count++
	s.sumExpectedGenerated += multiplier * (multiplier / (s.maxClientGenerate)) / 2
	resp := map[string]float64{
		"result": multiplier,
	}
	s.logger.Info("Values", zap.Float64("count", s.count), zap.Float64("p_zero", p_zero), zap.Float64("expected_RES", expectedRES))
	fmt.Println("Expected RES", expectedRES)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

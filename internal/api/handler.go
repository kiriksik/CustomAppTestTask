package api

import (
	"encoding/json"
	"fmt"
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
	alpha                float64
}

func NewServer(logger *zap.Logger, rtp float64) *Server {
	return &Server{
		logger:               logger,
		rtp:                  rtp,
		maxGenerate:          10000,
		maxClientGenerate:    10000,
		count:                0,
		sumExpectedGenerated: 0,
		alpha:                0.05,
	}
}

func (s *Server) RegisterRoutes() *http.ServeMux {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /get", s.HandleGet)
	return serveMux
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	/* expectedRES – сглаженное ожидаемое значение метрики RES

	alpha – коэффициент коррекции, который уменьшает expectedRES на каждой итерации,
	чтобы сократить разрыв между теоретическим и реальным значением.

	алгоритм работает так:
	рассчитывается мультипликатор (число от 0 до maxGenerate) путём разности между sumExpectedGenerated и суммой, которая должна быть в реальности.

	если expectedRES выше целевого RTP, p_zero = 1 (мультипликатор обнуляется),
	иначе p_zero = 0 (мультипликатор сохраняется).
	основное отличие в том что alpha компенсирует накопленное смещение,
	при большом числе итераций разница между expectedRES и реальным RES
	постепенно сокращается и алгоритм остаётся стабильным
	в том числе если все числа в последовательности одинаковые (исключается вероятность, которую даёт случайное генерирование мультипликатора
	и в первых шагах и в случае необходимости сильного изменения суммы мультипликатор = 10 000, т.е. любое число на стороне клиента в любом случае
	не будет пропущено */

	var multiplier float64
	var p_zero float64
	var expectedRES float64
	if s.count == 0 {
		expectedRES = 0
	} else {
		expectedRES = s.sumExpectedGenerated / s.count
		expectedRES = (1-s.alpha)*expectedRES + s.alpha
	}
	if expectedRES > s.rtp {
		p_zero = 1
	} else {
		p_zero = 0
		// meanGenerator := s.maxGenerate / 2
		// p_zero = 1 - (((s.count+1)*s.rtp - s.sumExpectedGenerated) / meanGenerator)
	}

	multiplier = (s.count+1)*s.rtp - s.sumExpectedGenerated
	if multiplier < 0 {
		multiplier = 0
	} else if multiplier > s.maxGenerate {
		multiplier = s.maxGenerate
	}
	if p_zero == 1 {
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

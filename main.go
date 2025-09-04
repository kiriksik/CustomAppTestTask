package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/kiriksik/CustomAppTestTask/internal/api"
	"github.com/kiriksik/CustomAppTestTask/internal/logger"
	"go.uber.org/zap"
)

var (
	rtpValue float64
)

func main() {
	// RTP
	rtpFlag := flag.Float64("rtp", 1.0, "target RTP value (0 < rtp <= 1.0)")
	flag.Parse()
	rtpValue = *rtpFlag

	httpPort := getEnv("PORT", "64333")

	// Logger
	logger.Init()
	defer logger.Sync()
	log := logger.L()
	log.Info("Запуск сервиса",
		zap.String("port", httpPort),
		zap.Float64("rtp", rtpValue),
	)

	// HTTP
	server := api.NewServer(log, rtpValue)
	mux := server.RegisterRoutes()
	if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
		log.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

package main

import (
	handler "main.go/handlers"
	"main.go/repository"
	"time"
)

func startDailyExpirationCheck(repo *repository.SystematicExamRepository, h *handler.SystematicExamHandler) {
	for {
		now := time.Now()
		next := now.Truncate(24 * time.Hour).Add(24 * time.Hour)
		time.Sleep(next.Sub(now))

		h.CheckExpiredExams()
	}
}

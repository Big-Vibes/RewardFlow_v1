package controller

import (
	"encoding/json"
	"net/http"
	"rewardpage/service"
	"rewardpage/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := utils.CreateContext()
	defer cancel()

	err := service.GetDB().Client().Ping(ctx, nil)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "db_error", "error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

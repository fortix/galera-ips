package server

import (
	"encoding/json"
	"net/http"

	"github.com/fortix/galera-ips/internal/monitor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Response struct {
	Writers []string `json:"writers"`
	Readers []string `json:"readers"`
}

func handleIps(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Writers: *monitor.Writers,
		Readers: *monitor.Readers,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func Run() {
	log.Info().Msgf("Listening on %s", viper.GetString("listen"))

	http.HandleFunc("/ips", handleIps)
	http.ListenAndServe(viper.GetString("listen"), nil)
}

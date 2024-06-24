package controller

import (
	"encoding/json"
	"fmt"
	"lesotho/global"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ApiKeyRequestBody struct {
	Client string `json:"client"`
}

type ApiKeyResponseBody struct {
	Key string `json:"key"`
}

func RequestApiKey(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msg(r.URL.EscapedPath())
	if r.Method != http.MethodPost {
		log.Error().Err(fmt.Errorf("method %s not allowed on %s", r.Method, r.URL.EscapedPath())).Send()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestBody ApiKeyRequestBody
	err := decoder.Decode(&requestBody)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apikey, err := global.ApiKeyRepo.IssueAPIKey(requestBody.Client)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info().Msgf("Issued an API key for client %s", requestBody.Client)

	result := ApiKeyResponseBody{Key: apikey}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

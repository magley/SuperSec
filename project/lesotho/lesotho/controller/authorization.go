package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"lesotho/acl"
	"lesotho/global"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

func AclUpdate(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msg(r.URL.EscapedPath())
	if r.Method != http.MethodPost {
		log.Error().Err(fmt.Errorf("method %s not allowed on %s", r.Method, r.URL.EscapedPath())).Send()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !global.ApiKeyRepo.CheckAPIKey(r) {
		log.Error().Err(fmt.Errorf("unauthorized request on %s", r.URL.EscapedPath())).Send()
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var aclDirective acl.ACLDirective
	err := decoder.Decode(&aclDirective)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = aclDirective.Validate()
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = global.Acl.Add(aclDirective, global.Nss)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info().Msgf("Added %v to the ACL.\n", aclDirective)
}

func AclQuery(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msgf(r.URL.EscapedPath())
	if r.Method != http.MethodGet {
		log.Error().Err(fmt.Errorf("method %s not allowed on %s", r.Method, r.URL.EscapedPath())).Send()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !global.ApiKeyRepo.CheckAPIKey(r) {
		log.Error().Err(fmt.Errorf("unauthorized request on %s", r.URL.EscapedPath())).Send()
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	aclDirective, err := acl.NewACLDirective(
		r.URL.Query().Get("object"),
		r.URL.Query().Get("relation"),
		r.URL.Query().Get("user"),
	)

	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorized := global.Acl.Check(aclDirective, global.Nss)
	result := AuthorizationResponse{Authorized: authorized}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func NamespaceUpdate(w http.ResponseWriter, r *http.Request) {
	log.Trace().Msgf(r.URL.EscapedPath())
	if r.Method != http.MethodPost {
		log.Error().Err(fmt.Errorf("method %s not allowed on %s", r.Method, r.URL.EscapedPath())).Send()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !global.ApiKeyRepo.CheckAPIKey(r) {
		log.Error().Err(fmt.Errorf("unauthorized request on %s", r.URL.EscapedPath())).Send()
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	namespaceAsString := new(strings.Builder)
	_, err := io.Copy(namespaceAsString, r.Body)
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace, err := global.Nss.Add(namespaceAsString.String())
	if err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Info().Msgf("Added Namespace %s", namespace.Name)
}

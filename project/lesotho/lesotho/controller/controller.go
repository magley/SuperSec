package controller

import (
	"encoding/json"
	"io"
	"lesotho/acl"
	"lesotho/global"
	ns "lesotho/namespace"
	"log"
	"net/http"
	"strings"
)

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

func AclUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\n", r.URL.EscapedPath())
	if r.Method != http.MethodPost {
		log.Printf("Method %s not allowed on %s", r.Method, r.URL.EscapedPath())
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var aclDirective acl.ACLDirective
	err := decoder.Decode(&aclDirective)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = aclDirective.Validate()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = global.Acl.Add(aclDirective, global.Nss)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Added %v to the ACL.\n", aclDirective)
}

func AclQuery(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\n", r.URL.EscapedPath())
	if r.Method != http.MethodGet {
		log.Printf("Method %s not allowed on %s", r.Method, r.URL.EscapedPath())
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	aclDirective, err := acl.NewACLDirective(
		r.URL.Query().Get("object"),
		r.URL.Query().Get("relation"),
		r.URL.Query().Get("user"),
	)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorized := global.Acl.Check(aclDirective, global.Nss)
	result := AuthorizationResponse{Authorized: authorized}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func NamespaceUpdate(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\n", r.URL.EscapedPath())
	if r.Method != http.MethodPost {
		log.Printf("Method %s not allowed on %s", r.Method, r.URL.EscapedPath())
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	namespaceAsString := new(strings.Builder)
	_, err := io.Copy(namespaceAsString, r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var namespace ns.Namespace
	err = json.Unmarshal([]byte(namespaceAsString.String()), &namespace)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	global.Nss.Add(namespace.Name, namespaceAsString.String())
	log.Printf("Added Namespace %s.\n", namespace.Name)
	log.Println(namespaceAsString.String())
}

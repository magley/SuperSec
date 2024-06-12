package main

import (
	"encoding/json"
	"lesotho/acl"
	ns "lesotho/namespace"
	"log"
	"net/http"
)

var glo_nsGraphCache *ns.NamespaceGraphCache
var glo_acl *acl.ACL
var glo_nss *ns.NamespaceStore

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

func aclUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	var aclDirective acl.ACLDirective
	err := decoder.Decode(&aclDirective)
	if err != nil {
		panic(err)
	}

	// TODO: Validation
	log.Println("TODO: Validate ACL Directive at POST /acl")

	glo_acl.Add(aclDirective)
	log.Printf("Added %v to the ACL.\n", aclDirective)
}

func aclQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	aclDirective := acl.NewACLDirective(
		r.URL.Query().Get("object"),
		r.URL.Query().Get("relation"),
		r.URL.Query().Get("user"),
	)

	authorized := glo_acl.Check(aclDirective, glo_nss)
	result := AuthorizationResponse{Authorized: authorized}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func namespaceUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	glo_nsGraphCache = ns.NewNamespaceGraphCache()

	glo_nss = ns.NewNamespaceStore(glo_nsGraphCache)
	glo_nss.AddFromFile("basic", "./basic.json")

	glo_acl = acl.NewACL("./data/acl/")
	glo_acl.AddFromFile("./basic.acl")
	defer glo_acl.Close()

	http.HandleFunc("/acl", aclUpdate)
	http.HandleFunc("/acl/check", aclQuery)
	http.HandleFunc("/namespace", namespaceUpdate)

	log.Println("Serving http://127.0.0.1:5000")

	http.ListenAndServe("127.0.0.1:5000", nil)
}

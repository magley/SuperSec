package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var glo_acl *ACL
var glo_nss *NamespaceStore

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

func aclUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	var t ACLDirective
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	glo_acl.Add(t.Object, t.Relation, t.User)
	log.Printf("Added %s#%s@%s to the ACL.\n", t.Object, t.Relation, t.User)
}

func aclQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	t := ACLDirective{
		Object:   r.URL.Query().Get("object"),
		Relation: r.URL.Query().Get("relation"),
		User:     r.URL.Query().Get("user"),
	}

	authorized := glo_acl.Check(t.Object, t.Relation, t.User, glo_nss)
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
	glo_nss = NewNamespaceStore()
	glo_nss.AddFromFile("basic", "./basic.json")

	glo_acl = NewACL("./data/acl/")
	glo_acl.AddFromFile("./basic.acl")

	http.HandleFunc("/acl", aclUpdate)
	http.HandleFunc("/acl/check", aclQuery)
	http.HandleFunc("/namespace", namespaceUpdate)

	log.Println("Serving http://127.0.0.1:5000")

	http.ListenAndServe("127.0.0.1:5000", nil)
}

package main

import (
	"encoding/json"
	"io"
	"lesotho/acl"
	ns "lesotho/namespace"
	"log"
	"net/http"
	"strings"
)

var glo_acl *acl.ACL
var glo_nss *ns.NamespaceStore

type AuthorizationResponse struct {
	Authorized bool `json:"authorized"`
}

func aclUpdate(w http.ResponseWriter, r *http.Request) {
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

	err = glo_acl.Add(aclDirective, glo_nss)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Added %v to the ACL.\n", aclDirective)
}

func aclQuery(w http.ResponseWriter, r *http.Request) {
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

	authorized := glo_acl.Check(aclDirective, glo_nss)
	result := AuthorizationResponse{Authorized: authorized}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func namespaceUpdate(w http.ResponseWriter, r *http.Request) {
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

	glo_nss.Add(namespace.Name, namespaceAsString.String())
	log.Printf("Added Namespace %s.\n", namespace.Name)
	log.Println(namespaceAsString.String())
}

func main() {
	namespaceGraphCache := ns.NewNamespaceGraphCache()

	glo_nss = ns.NewNamespaceStore(namespaceGraphCache)
	glo_nss.AddFromFile("basic", "./basic.json")

	glo_acl = acl.NewACL("./data/acl/")
	glo_acl.AddFromFile("./basic.acl", glo_nss)
	defer glo_acl.Close()

	http.HandleFunc("/acl", aclUpdate)
	http.HandleFunc("/acl/check", aclQuery)
	http.HandleFunc("/namespace", namespaceUpdate)

	log.Println("Serving http://127.0.0.1:5000")

	http.ListenAndServe("127.0.0.1:5000", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lesotho/acl"
	ns "lesotho/namespace"
	"log"
	"net/http"
	"strings"

	"gopkg.in/ini.v1"
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
	log.Println("Loading configuration from config.ini ...")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Fail to read 'config.ini': %v", err)
		return
	}

	log.Println("Parsing configuration  ...")
	cfg_ip := cfg.Section("MAIN").Key("ip").String()
	cfg_port := cfg.Section("MAIN").Key("port").String()
	cfg_ns_fname := ""
	cfg_ns_name := ""
	cfg_acl_path := cfg.Section("ACL").Key("path").String()
	cfg_acl_fname := ""
	cfg_use_cache := cfg.Section("MAIN").Key("use_graph_namespace_cache").MustBool(true)

	k := cfg.Section("NAMESPACE").Key("namespace")
	if k != nil {
		cfg_ns_fname = k.String()
	}
	k = cfg.Section("NAMESPACE").Key("namespace_name")
	if k != nil {
		cfg_ns_name = k.String()
	}
	k = cfg.Section("ACL").Key("acl")
	if k != nil {
		cfg_acl_fname = k.String()
	}

	var namespaceGraphCache *ns.NamespaceGraphCache
	if cfg_use_cache {
		log.Println("Building namespace graph cache ...")
		namespaceGraphCache = nil
	} else {
		log.Println("Namespace graph cache is ignored, skipping ...")
	}

	log.Println("Building namespace store ...")
	glo_nss = ns.NewNamespaceStore(namespaceGraphCache)

	if cfg_ns_fname != "" {
		log.Printf("Loading namespace '%s' from '%s' ...\n", cfg_ns_name, cfg_ns_fname)
		glo_nss.AddFromFile(cfg_ns_name, cfg_ns_fname)
	}

	glo_acl = acl.NewACL(cfg_acl_path)
	if cfg_acl_fname != "" {
		log.Printf("Loading ACL from '%s' ...\n", cfg_acl_fname)
		glo_acl.AddFromFile(cfg_acl_fname, glo_nss)
	}
	defer glo_acl.Close()

	http.HandleFunc("/acl", aclUpdate)
	http.HandleFunc("/acl/check", aclQuery)
	http.HandleFunc("/namespace", namespaceUpdate)

	lesotho_host := fmt.Sprintf("%s:%s", cfg_ip, cfg_port)
	log.Printf("Serving Lesotho on http://%s\n", lesotho_host)
	http.ListenAndServe(lesotho_host, nil)
}

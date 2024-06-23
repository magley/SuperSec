package main

import (
	"fmt"
	"lesotho/acl"
	"lesotho/controller"
	"lesotho/global"
	ns "lesotho/namespace"
	"log"
	"net/http"

	"gopkg.in/ini.v1"
)

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
	global.Nss = ns.NewNamespaceStore(namespaceGraphCache)

	if cfg_ns_fname != "" {
		log.Printf("Loading namespace '%s' from '%s' ...\n", cfg_ns_name, cfg_ns_fname)
		global.Nss.AddFromFile(cfg_ns_name, cfg_ns_fname)
	}

	global.Acl = acl.NewACL(cfg_acl_path)
	if cfg_acl_fname != "" {
		log.Printf("Loading ACL from '%s' ...\n", cfg_acl_fname)
		global.Acl.AddFromFile(cfg_acl_fname, global.Nss)
	}
	defer global.Acl.Close()

	http.HandleFunc("/acl", controller.AclUpdate)
	http.HandleFunc("/acl/check", controller.AclQuery)
	http.HandleFunc("/namespace", controller.NamespaceUpdate)

	lesotho_host := fmt.Sprintf("%s:%s", cfg_ip, cfg_port)
	log.Printf("Serving Lesotho on http://%s\n", lesotho_host)
	http.ListenAndServe(lesotho_host, nil)
}

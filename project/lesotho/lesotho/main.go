package main

import (
	"fmt"
	"lesotho/acl"
	"lesotho/apikey"
	"lesotho/controller"
	"lesotho/global"
	ns "lesotho/namespace"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	fileWriter := &lumberjack.Logger{
		Filename: "./logs/lesotho.log",
		MaxSize:  1,
	}

	multiLevelWriter := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	log.Logger = zerolog.New(multiLevelWriter).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
}

func main() {
	initLogger()

	log.Info().Msgf("Loading configuration from config.ini ...")
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Error().Err(err).Msg("Fail to read 'config.ini'")
		return
	}

	log.Info().Msgf("Parsing configuration  ...")
	cfg_ip := cfg.Section("MAIN").Key("ip").String()
	cfg_port := cfg.Section("MAIN").Key("port").String()
	cfg_ns_fname := ""
	cfg_api_key_repo_path := cfg.Section("MAIN").Key("api_key_repo_path").String()

	cfg_acl_path := cfg.Section("ACL").Key("path").String()
	cfg_acl_fname := ""
	cfg_use_cache := cfg.Section("MAIN").Key("use_graph_namespace_cache").MustBool(true)

	log.Info().Msgf("Loading API key repository ...")
	global.ApiKeyRepo = apikey.OpenAPIKeyRepository(cfg_api_key_repo_path)
	defer global.ApiKeyRepo.Close()

	k := cfg.Section("NAMESPACE").Key("namespace")
	if k != nil {
		cfg_ns_fname = k.String()
	}
	k = cfg.Section("ACL").Key("acl")
	if k != nil {
		cfg_acl_fname = k.String()
	}

	var namespaceGraphCache *ns.NamespaceGraphCache
	if cfg_use_cache {
		log.Info().Msgf("Building namespace graph cache ...")
		namespaceGraphCache = ns.NewNamespaceGraphCache()
	} else {
		log.Info().Msgf("Namespace graph cache is ignored, skipping ...")
	}

	log.Info().Msgf("Building namespace store ...")
	global.Nss = ns.NewNamespaceStore(namespaceGraphCache)

	if cfg_ns_fname != "" {
		log.Info().Msgf("Loading namespace from '%s' ...", cfg_ns_fname)
		ns, err := global.Nss.AddFromFile(cfg_ns_fname)
		if err != nil {
			log.Error().Err(err).Msgf("Failed loading namespace from '%s'", cfg_ns_fname)
			return
		} else {
			log.Info().Msgf("Succeeded loading namespace '%s' from '%s'", ns.Name, cfg_ns_fname)
		}
	}

	global.Acl = acl.NewACL(cfg_acl_path)
	if cfg_acl_fname != "" {
		log.Info().Msgf("Loading ACL from '%s' ...", cfg_acl_fname)
		global.Acl.AddFromFile(cfg_acl_fname, global.Nss)
	}
	defer global.Acl.Close()

	log.Info().Msgf("Loading HTTPs configuration ... ")
	cfg_https_use, err := cfg.Section("HTTPS").Key("use_https").Bool()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}
	cfg_https_crt := ""
	cfg_https_key := ""

	if cfg_https_use {
		cfg_https_crt = cfg.Section("HTTPS").Key("https_cert").String()
		cfg_https_key = cfg.Section("HTTPS").Key("https_key").String()
	}

	http.HandleFunc("/acl", controller.AclUpdate)
	http.HandleFunc("/acl/check", controller.AclQuery)
	http.HandleFunc("/namespace", controller.NamespaceUpdate)
	http.HandleFunc("/apikey", controller.RequestApiKey)

	lesotho_host := fmt.Sprintf("%s:%s", cfg_ip, cfg_port)

	if cfg_https_use {
		log.Info().Msgf("Serving Lesotho on https://%s", lesotho_host)
		http.ListenAndServeTLS(lesotho_host, cfg_https_crt, cfg_https_key, nil)
	} else {
		log.Info().Msgf("Serving Lesotho on http://%s", lesotho_host)
		http.ListenAndServe(lesotho_host, nil)
	}
}

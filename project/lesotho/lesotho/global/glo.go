package global

import (
	"lesotho/acl"
	"lesotho/apikey"
	ns "lesotho/namespace"
)

var Acl *acl.ACL
var Nss *ns.NamespaceStore
var ApiKeyRepo *apikey.APIKeyRepository
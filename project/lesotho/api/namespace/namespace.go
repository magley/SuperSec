package namespace

import (
	"encoding/json"
	"os"

	capi "github.com/hashicorp/consul/api"
)

type NamespaceStore struct {
	client *capi.Client
}

type NamespaceRelationUnionElementUserset struct {
	Relation string `json:"relation"`
}

type NamespaceRelationUnionElement struct {
	Self            *map[string]interface{}               `json:"self,omitempty"`
	ComputedUserset *NamespaceRelationUnionElementUserset `json:"computed_userset,omitempty"`
}

type NamespaceRelation struct {
	Union *[]NamespaceRelationUnionElement `json:"union,omitempty"`
}

type Namespace struct {
	Name      string                       `json:"name"`
	Relations map[string]NamespaceRelation `json:"relations"`
}

func NewNamespaceStore() *NamespaceStore {
	nss := new(NamespaceStore)
	return nss
}

//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/
// Low level methods.
// Abstracting ConsulDB through NamespaceStore.
//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/

func (nss *NamespaceStore) openClientIfNeeded() {
	if nss.client == nil {
		client, err := capi.NewClient(capi.DefaultConfig())
		if err != nil {
			panic(err)
		}
		nss.client = client
	}
}

func (nss *NamespaceStore) Get(key string) string {
	nss.openClientIfNeeded()
	kv := nss.client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		panic(err)
	}
	return string(pair.Value)
}

func (nss *NamespaceStore) Put(key string, val string) {
	nss.openClientIfNeeded()
	kv := nss.client.KV()

	p := &capi.KVPair{Key: key, Value: []byte(val)}
	_, err := kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
}

//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/
// High level methods.
//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/

func (nss *NamespaceStore) Add(namespaceName string, namespacJson string) {
	nss.Put(namespaceName, namespacJson)
}

func (nss *NamespaceStore) AddFromFile(namespaceName string, namespaceDataFname string) {
	data, err := os.ReadFile(namespaceDataFname)
	if err != nil {
		panic(err)
	}
	nss.Add(namespaceName, string(data))
}

func (nss *NamespaceStore) GetRelations(namespaceName string) map[string]NamespaceRelation {
	namespaceJson := nss.Get(namespaceName)

	var namespace Namespace
	err := json.Unmarshal([]byte(namespaceJson), &namespace)
	if err != nil {
		panic(err)
	}

	relations := namespace.Relations

	return relations
}

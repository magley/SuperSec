package namespace

import (
	"encoding/json"
	"os"

	capi "github.com/hashicorp/consul/api"
)

type NamespaceStore struct {
	client     *capi.Client
	graphCache *NamespaceGraphCache
}

func NewNamespaceStore(graphCache *NamespaceGraphCache) *NamespaceStore {
	nss := new(NamespaceStore)
	nss.graphCache = graphCache
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

func (nss *NamespaceStore) Add(namespaceName string, namespaceJson string) {
	nss.Put(namespaceName, namespaceJson)
	nss.revalidateGraphCache(namespaceName, namespaceJson)
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
	return nss.GetRelationsFromJson(namespaceJson)
}

func (nss *NamespaceStore) GetRelationsFromJson(namespaceJson string) map[string]NamespaceRelation {
	var namespace Namespace
	err := json.Unmarshal([]byte(namespaceJson), &namespace)
	if err != nil {
		panic(err)
	}

	relations := namespace.Relations

	return relations
}

func (nss *NamespaceStore) revalidateGraphCache(namespaceName string, namespaceJson string) {
	if nss.graphCache == nil {
		return
	}

	graph, ok := nss.graphCache.Get(namespaceName)
	if ok {
		relations := nss.GetRelationsFromJson(namespaceJson)
		graph.RebuildFromNamespaceRelations(relations)
		nss.graphCache.Put(namespaceName, graph)
	}
}

func (nss *NamespaceStore) GetNamespaceGraph(namespaceName string) *NamespaceGraph {
	if nss.graphCache == nil {
		return nss.buildGraph(namespaceName)
	} else {
		graph, ok := nss.graphCache.Get(namespaceName)
		if !ok {
			graph = nss.buildGraph(namespaceName)
			nss.graphCache.Put(namespaceName, graph)
		}
		return graph
	}
}

func (nss *NamespaceStore) buildGraph(namespaceName string) *NamespaceGraph {
	graph := NewNamespaceGraph()
	relations := nss.GetRelations(namespaceName)
	graph.RebuildFromNamespaceRelations(relations)
	return graph
}

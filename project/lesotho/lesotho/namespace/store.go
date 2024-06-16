package namespace

import (
	"encoding/json"
	"fmt"
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

func (nss *NamespaceStore) Get(key string) (string, error) {
	nss.openClientIfNeeded()
	kv := nss.client.KV()

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return "", err
	}
	if pair == nil {
		return "", fmt.Errorf("namespace %s not found in namespace store", key)
	}
	return string(pair.Value), nil
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

func (nss *NamespaceStore) HasNamespace(namespaceName string) (bool, error) {
	if nss.graphCache == nil {
		_, err := nss.Get(namespaceName)
		ok := err == nil
		return ok, err
	} else {
		_, ok := nss.graphCache.Get(namespaceName)
		if !ok {
			_, err := nss.Get(namespaceName)
			ok := err == nil
			return ok, err
		}
		return true, nil
	}
}

func (nss *NamespaceStore) GetRelations(namespaceName string) (map[string]NamespaceRelation, error) {
	namespaceJson, err := nss.Get(namespaceName)
	if err != nil {
		return nil, err
	}
	return nss.GetRelationsFromJson(namespaceJson), nil
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

func (nss *NamespaceStore) GetNamespaceGraph(namespaceName string) (*NamespaceGraph, error) {
	if nss.graphCache == nil {
		return nss.buildGraph(namespaceName)
	} else {
		graph, ok := nss.graphCache.Get(namespaceName)
		if !ok {
			graph_, err := nss.buildGraph(namespaceName)
			if err != nil {
				return nil, err
			}
			graph = graph_ // This is neccessary because without graph_, the
			// regular graph gets redeclared since it's in an inner scope.
			nss.graphCache.Put(namespaceName, graph)
		}
		return graph, nil
	}
}

func (nss *NamespaceStore) buildGraph(namespaceName string) (*NamespaceGraph, error) {
	graph := NewNamespaceGraph()
	relations, err := nss.GetRelations(namespaceName)
	if err != nil {
		return nil, err
	}

	graph.RebuildFromNamespaceRelations(relations)
	return graph, nil
}

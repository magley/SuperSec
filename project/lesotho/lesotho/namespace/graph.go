package namespace

import (
	"fmt"
	"slices"
)

// NamespaceGraph is a graph data structure describing a namespace object.
type NamespaceGraph struct {
	// Relations is a JSON-like map representing the namespace.
	//
	// Key is the name of the relations e.g. "owner", "reviewer", "editor".
	//
	// Value is a list of parent relations i.e. relations that "inherit" all
	// permissions from this relation. It probably makes more sense to call them
	// child relations, though.
	Relations map[string][]string
}

func NewNamespaceGraph() *NamespaceGraph {
	g := new(NamespaceGraph)
	g.Relations = make(map[string][]string)
	return g
}

func NamespaceGraphFromNamespaceRelations(relations map[string]NamespaceRelation) *NamespaceGraph {
	g := NewNamespaceGraph()
	for relName, relContent := range relations {
		g.Relations[relName] = make([]string, 0)

		if relContent.Union != nil {
			unionElements := *relContent.Union
			for _, unionElement := range unionElements {
				if unionElement.ComputedUserset != nil {
					parent := unionElement.ComputedUserset.Relation
					g.Relations[relName] = append(g.Relations[relName], parent)
				}
			}
		}
	}
	return g
}

func (g *NamespaceGraph) ValidateNamespaceGraph() error {
	for rel, parents := range g.Relations {
		for _, p := range parents {
			// Rule 1
			if rel == p {
				return fmt.Errorf("relation '%s' has itself in computed_userset", rel)
			}

			// Rule 2
			if slices.Contains(g.Relations[p], rel) {
				return fmt.Errorf("relations '%s' and '%s' are in each other's computed_userset (circular reference)", rel, p)
			}
		}
	}
	return nil
}

package namespace

import "fmt"

// NamespaceGraph is a graph data structure describing a namespace object.
type NamespaceGraph struct {
	// Relations is a JSON-like map representing the namespace.
	//
	// Key is the name of the relations e.g. "owner", "reviewer", "editor".
	//
	// Value is a list of child relations i.e. relations that "inherit" all
	// permissions from this relation.
	Relations map[string][]string
}

func NewNamespaceGraph() *NamespaceGraph {
	g := new(NamespaceGraph)
	g.Relations = make(map[string][]string)
	return g
}

func (g *NamespaceGraph) RebuildFromNamespaceRelations(relations map[string]NamespaceRelation) {
	for relName, relContent := range relations {
		g.Relations[relName] = make([]string, 0)

		if relContent.Union != nil {
			unionElements := *relContent.Union
			for _, unionElement := range unionElements {
				if unionElement.ComputedUserset != nil {
					child := unionElement.ComputedUserset.Relation
					g.Relations[relName] = append(g.Relations[relName], child)
				}
			}
		}

		fmt.Printf("%s %+v\n", relName, relContent)
	}
}

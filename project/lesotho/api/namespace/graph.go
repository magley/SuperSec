package namespace

type NamespaceGraph struct {
	G map[string][]string
}

func NewNamespaceGraph() *NamespaceGraph {
	g := new(NamespaceGraph)
	g.G = make(map[string][]string)
	return g
}

func (g *NamespaceGraph) RebuildFromNamespaceRelations(relations map[string]NamespaceRelation) {
	for relName, relContent := range relations {
		g.G[relName] = make([]string, 0)

		if relContent.Union != nil {
			unionElements := *relContent.Union
			for _, unionElement := range unionElements {
				if unionElement.ComputedUserset != nil {
					child := unionElement.ComputedUserset.Relation
					g.G[relName] = append(g.G[relName], child)
				}
			}
		}
	}
}

package namespace

type Namespace struct {
	Name      string                       `json:"name"`
	Relations map[string]NamespaceRelation `json:"relations"`
}

type NamespaceRelation struct {
	Union *[]NamespaceRelationUnionElement `json:"union,omitempty"`
}

type NamespaceRelationUnionElement struct {
	Self            *map[string]interface{}               `json:"self,omitempty"`
	ComputedUserset *NamespaceRelationUnionElementUserset `json:"computed_userset,omitempty"`
}

type NamespaceRelationUnionElementUserset struct {
	Relation string `json:"relation"`
}

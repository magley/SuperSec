package acl

import (
	"fmt"
	"strings"
)

// ACLDirective is a triplet `object#relation@user` used in access control
// lists.
type ACLDirective struct {
	Object   string
	Relation string
	User     string
}

// Validate checks if any field in the directive contains illegal characters.
func (d *ACLDirective) Validate() error {
	objectParts := strings.Split(d.Object, ":")

	if len(objectParts) != 2 {
		return fmt.Errorf("field 'object' (%s) in ACLDirective (%s) has invalid format", d.Object, d.String())
	}

	justNamespace := objectParts[0]
	justObject := objectParts[1]

	if justNamespace == "" {
		return fmt.Errorf("field object in ACLDirective (%s) has empty namespace", d.String())
	}
	if strings.ContainsAny(justNamespace, ":#@") {
		return fmt.Errorf("field object (%s) in ACLDirective (%s) contains invalid character", d.Object, d.String())
	}
	if justObject == "" {
		return fmt.Errorf("field object in ACLDirective (%s) is an empty string", d.String())
	}
	if strings.ContainsAny(justObject, ":#@") {
		return fmt.Errorf("field object (%s) in ACLDirective (%s) contains invalid character", d.Object, d.String())
	}

	if d.Relation == "" {
		return fmt.Errorf("field relation in ACLDirective (%s) is an empty string", d.String())
	}
	if strings.ContainsAny(d.Relation, ":#@") {
		return fmt.Errorf("field relation (%s) in ACLDirective (%s) contains invalid character", d.Relation, d.String())
	}

	if d.User == "" {
		return fmt.Errorf("field user in ACLDirective (%s) is an empty string", d.String())
	}
	if strings.ContainsAny(d.User, ":#@") {
		return fmt.Errorf("field user (%s) in ACLDirective (%s) contains invalid character", d.User, d.String())
	}
	return nil
}

// Assumes ACLDirective is valid
func (d *ACLDirective) Namespace() string {
	return strings.Split(d.Object, ":")[0]
}

// String converts the ACLDirective into canonical form.
func (d *ACLDirective) String() string {
	return fmt.Sprintf("%s#%s@%s", d.Object, d.Relation, d.User)
}

func (d *ACLDirective) ObjectUserString() string {
	return fmt.Sprintf("%s-%s", d.Object, d.User)
}

func NewACLDirective(object string, relation string, user string) (*ACLDirective, error) {
	aclDirective := new(ACLDirective)

	aclDirective.Object = object
	aclDirective.Relation = relation
	aclDirective.User = user

	err := aclDirective.Validate()
	if err != nil {
		return nil, err
	}

	return aclDirective, nil
}

func newACLDirectiveWithoutValidation(object string, relation string, user string) *ACLDirective {
	aclDirective := new(ACLDirective)

	aclDirective.Object = object
	aclDirective.Relation = relation
	aclDirective.User = user

	return aclDirective
}

func NewACLDirectiveFromCanonicalString(canonical string) (*ACLDirective, error) {
	parts := strings.Split(canonical, "#")
	if len(parts) != 2 {
		return nil, fmt.Errorf("canonical string (%s) has invalid format", canonical)
	}
	object := parts[0]
	parts = strings.Split(parts[1], "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("canonical string (%s) has invalid format", canonical)
	}
	relation := parts[0]
	user := parts[1]
	return NewACLDirective(object, relation, user)
}

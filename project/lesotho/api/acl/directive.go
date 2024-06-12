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
		return fmt.Errorf("field object (%s) in ACLDirective has invalid format", d.Object)
	}

	justNamespace := objectParts[0]
	justObject := objectParts[1]

	if strings.ContainsAny(justNamespace, ":#@") {
		return fmt.Errorf("field object (%s) in ACLDirective contains invalid character", d.Object)
	}
	if strings.ContainsAny(justObject, ":#@") {
		return fmt.Errorf("field object (%s) in ACLDirective contains invalid character", d.Object)
	}
	if strings.ContainsAny(d.Relation, ":#@") {
		return fmt.Errorf("field relation (%s) in ACLDirective contains invalid character", d.Relation)
	}
	if strings.ContainsAny(d.User, ":#@") {
		return fmt.Errorf("field user (%s) in ACLDirective contains invalid character", d.User)
	}
	return nil
}

// String converts the ACLDirective into canonical form.
func (d *ACLDirective) String() string {
	return fmt.Sprintf("%s#%s@%s", d.Object, d.Relation, d.User)
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

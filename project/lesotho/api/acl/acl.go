package acl

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"

	ns "lesotho/namespace"
)

type ACLDirective struct {
	Object   string
	Relation string
	User     string
}

func (d *ACLDirective) String() string {
	return fmt.Sprintf("%s#%s@%s", d.Object, d.Relation, d.User)
}

func NewACLDirective(object string, relation string, user string) *ACLDirective {
	aclDirective := new(ACLDirective)

	// TODO: Validation

	aclDirective.Object = object
	aclDirective.Relation = relation
	aclDirective.User = user

	return aclDirective
}

type ACL struct {
	fname string
	db    *leveldb.DB
}

func NewACL(fname string) *ACL {
	acl := new(ACL)
	acl.fname = fname

	db, err := leveldb.OpenFile(acl.fname, nil)
	if err != nil {
		panic(err)
	}
	acl.db = db
	return acl
}

func (acl *ACL) Close() {
	if acl.db != nil {
		acl.db.Close()
	}
}

//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/
// Low level methods.
// Abstracting LevelDB through ACL.
//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/

func (acl *ACL) Get(key string) string {
	val, err := acl.db.Get([]byte(key), nil)
	if err != nil {
		panic(err)
	}

	return string(val)
}

func (acl *ACL) Put(key string) {
	err := acl.db.Put([]byte(key), []byte{}, nil)
	if err != nil {
		panic(err)
	}
}

func (acl *ACL) Has(key string) bool {
	has, err := acl.db.Has([]byte(key), nil)
	if err != nil {
		panic(err)
	}

	return has
}

//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/
// High level methods.
//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/

func (acl *ACL) addDirective(directive string) {
	acl.Put(directive)
}

func (acl *ACL) Add(aclDirective ACLDirective) {
	acl.addDirective(aclDirective.String())
}

func (acl *ACL) AddFromFile(aclDataFname string) {
	file, err := os.Open(aclDataFname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		acl.addDirective(line)
	}
}

func removeDuplicates(li []string) []string {
	slices.Sort(li)
	return slices.Compact(li)
}

func (acl *ACL) Check(aclDirective *ACLDirective, nss *ns.NamespaceStore) bool {
	directive := aclDirective.String()

	if acl.Has(directive) {
		return true
	}

	parts := strings.Split(aclDirective.Object, ":")
	if len(parts) != 2 {
		panic("object in an ACL directive must have the following structure: name:instance")
	}
	namespaceName := parts[0]
	G := nss.GetNamespaceGraph(namespaceName)

	relationParents := make([]string, 0)
	queue := []string{aclDirective.Relation}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		// Get direct parent of e
		dp := make([]string, 0)
		for k, v := range G.Relations {
			if e == k {
				dp = append(dp, v...)
			}
		}

		relationParents = append(relationParents, dp...)
		queue = append(queue, dp...)

		queue = removeDuplicates(queue)
	}

	relationParents = removeDuplicates(relationParents)
	for _, r := range relationParents {
		aclD := NewACLDirective(aclDirective.Object, r, aclDirective.User)
		if acl.Has(aclD.String()) {
			return true
		}
	}

	return false
}

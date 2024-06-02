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

func (acl *ACL) Put(key string, val string) {
	err := acl.db.Put([]byte(key), []byte(val), nil)
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
	acl.Put(directive, "")
}

func (acl *ACL) Add(object string, relation string, user string) {
	directive := fmt.Sprintf("%s#%s@%s", object, relation, user)
	acl.addDirective(directive)
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

func (acl *ACL) Check(object string, relation string, user string, nss *ns.NamespaceStore) bool {
	directive := fmt.Sprintf("%s#%s@%s", object, relation, user)

	if acl.Has(directive) {
		return true
	}

	parts := strings.Split(object, ":")
	if len(parts) != 2 {
		panic("object in an ACL directive must have the following structure: name:instance")
	}
	namespaceName := parts[0]
	relations := nss.GetRelations(namespaceName)

	G := ns.NewNamespaceGraph()
	G.RebuildFromNamespaceRelations(relations)

	relationParents := make([]string, 0)
	queue := []string{relation}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		// Get direct parent of e
		dp := make([]string, 0)
		for k, v := range G.G {
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
		directive := fmt.Sprintf("%s#%s@%s", object, r, user)

		if acl.Has(directive) {
			return true
		}
	}

	return false
}
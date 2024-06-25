package acl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/syndtr/goleveldb/leveldb"

	ns "lesotho/namespace"
)

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

func (acl *ACL) Has(directive *ACLDirective) bool {
	val, err := acl.db.Get([]byte(directive.ObjectUserString()), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false
		} else {
			panic(err)
		}
	}

	return directive.Relation == string(val)
}

//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/
// High level methods.
//_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/_/

// ACLDirective is stored as key-val pair where key is {object}-{user} and val is {relation}
func (acl *ACL) Add(aclDirective ACLDirective, nss *ns.NamespaceStore) error {
	namespace := aclDirective.Namespace()

	_, err := nss.HasNamespace(namespace)
	if err != nil {
		return err
	}

	canBeSet, err := nss.IsSettable(namespace, aclDirective.Relation)
	if !canBeSet {
		return fmt.Errorf("relation '%s' cannot be set in namespace '%s' (tried adding ACL directive '%s')", aclDirective.Relation, namespace, aclDirective.String())
	}
	if err != nil {
		return err
	}

	acl.Put(aclDirective.ObjectUserString(), aclDirective.Relation)
	return nil
}

func (acl *ACL) AddFromFile(aclDataFname string, nss *ns.NamespaceStore) {
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

		directive, err := NewACLDirectiveFromCanonicalString(line)
		if err != nil {
			panic(err)
		}
		err = acl.Add(*directive, nss)
		if err != nil {
			panic(err)
		}
	}
}

func (acl *ACL) Check(aclDirective *ACLDirective, nss *ns.NamespaceStore) bool {
	if acl.Has(aclDirective) {
		return true
	}

	namespaceName := aclDirective.Namespace()
	G, err := nss.GetNamespaceGraph(namespaceName)
	if err != nil {
		log.Error().Err(err).Send()
		return false
	}

	relationParents := make([]string, 0)
	queue := []string{aclDirective.Relation}
	queueUsed := map[string]bool{}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		_, alreadyUsed := queueUsed[e]
		if alreadyUsed {
			continue
		}
		queueUsed[e] = true

		// Get direct parent of e
		dp := make([]string, 0)
		v, ok := G.Relations[e]
		if ok {
			dp = append(dp, v...)
		}

		relationParents = append(relationParents, dp...)
		queue = append(queue, dp...)
	}

	parentsChecked := map[string]bool{}
	for _, r := range relationParents {
		_, alreadyUsed := parentsChecked[r]
		if alreadyUsed {
			continue
		}
		parentsChecked[r] = true

		aclD := newACLDirectiveWithoutValidation(aclDirective.Object, r, aclDirective.User)
		if acl.Has(aclD) {
			return true
		}
	}

	return false
}

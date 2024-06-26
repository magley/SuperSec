# About

This is a prototype written in Python, using in-memory data.
The goal of this prototype is to experiment with the algorithm for determining
granted permissions.

# Model

It's written under the followiong assumptions:

- Namespaces are saved in a key value store, where the key is the name of the
  namespace (or some sort of namespace identifier) and value is the json string
  of the namespace (see `basic.json` for an example).
- ACLs are saved in a key value store, where the key is a single ACL directive
  and the value is "nothing" (see `basic.acl` for an example)
- Namespaces only support the `union` operator.

# Getting started

Run the server with `run.bat`. Test the `basic` namespace and ACL by running
`test.py`.

# Algorithm

Consider the pseudocode:

```py
def has_access(object, relation, user):
    True if 'object#relation@user' in ACL

    parents = get_parent_relations_in_namespace(relation, object.namespace)
    for parent_relation in parents:
         True if 'object#parent_relation@user' in ACL

    False
```

We need to build a tree (forest?) from the namespace's JSON and then traverse
through the entire graph. Simplified, the namespace looks like this:

```json
"relations": {
    "r1": [],
    "r2": [],
    "r3": ["r1"],
    "r4": ["r1", "r2"],
    "r5": ["r3"],
}
```

Each relation defines zero or more "parent" relations. Continuing with the above
example, if user `u` has `r1`, then `u` also has `r3` and `r4`. Therefore, when
querying `u` for `r3`, we need to do a breadth first search on its parents.

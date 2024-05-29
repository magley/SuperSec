import json

#   key = namespace "id"
#   value = json of namespace contents
glo_namespaces = {}

#   key = acl
#   value = nothing
glo_acl = {}


def add_namespace(namespace_name: str, namespace_json: str):
    global glo_namespaces

    glo_namespaces[namespace_name] = namespace_json


def _add_acl(directive: str):
    global glo_acl
    glo_acl[directive] = 0


def add_acl(object: str, relation: str, user: str):
    directive = f"{object}#{relation}@{user}"
    _add_acl(directive)


def add_acl_from_file(fname: str):
    with open(fname) as f:
        for line in f.readlines():
            line2 = line.strip()
            if line2.startswith("#"):
                continue
            if len(line2) == 0:
                continue

            _add_acl(line2)


def check_acl(object: str, relation: str, user: str) -> bool:
    # Check for permission directly.
    #
    directive = f"{object}#{relation}@{user}"
    if directive in glo_acl:
        return True

    # Bake graph from namespace.
    #
    parts = object.split(":")
    namespace_name = parts[0]

    namespace = json.loads(glo_namespaces[namespace_name])
    relations: dict = namespace["relations"]

    #   key = relation (string)
    #   value = list of child relations (string) for this key
    G = {}

    for rel_name, rel_content in relations.items():
        G[rel_name] = []

        if "union" in rel_content:
            union_elements: list = rel_content["union"]
            for union_element in union_elements:
                if "computed_userset" in union_element:
                    child = union_element["computed_userset"]["relation"]
                    G[rel_name].append(child)

    # Get all parents of `relation`
    #
    relation_parents = []

    def get_direct_parent_of(e):
        direct_parents = []
        for k, v in G.items():
            if e == k:
                direct_parents += v
        return direct_parents

    queue = [relation]
    while len(queue) > 0:
        e = queue[0]
        queue = queue[1:]

        dp = get_direct_parent_of(e)
        relation_parents += dp
        queue += dp

        queue = list(set(queue))

    relation_parents = set(relation_parents)
    print(relation_parents)
    for r in relation_parents:
        directive = f"{object}#{r}@{user}"
        if directive in glo_acl:
            return True

    return False

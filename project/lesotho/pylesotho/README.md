# Pylesotho

A Python library providing an easy interface to the Lesotho API.

## Installing (local)

```sh
pip install ./pylesotho
```

## Usage

```py
from pylesotho.client import LesothoClient

# Create a client.
# Make sure to obtain an API key.

cli = LesothoClient("http://localhost:5000", LESOTHO_CLIENT_NAME, LESOTHO_API_KEY)
cli = LesothoClient("https://localhost:5000", LESOTHO_CLIENT_NAME, LESOTHO_API_KEY, './server.cert') # If using HTTPs


# Upload a namespace called "namespace1"

namespace = {
    "name": "namespace1",
    "relations": {
        "owner": {},
        "editor": {
            "union": [
                { "this": {} },
                { "computed_userset": { "relation": "owner" } }
            ]
        },
        "viewer": {
            "union": [
                { "this": {} },
                { "computed_userset": { "relation": "editor" } }
            ]
        }
    }
}

cli.namespace_update(namespace)

# Insert directive to the access control list.

cli.acl_update("namespace1", "file1", "owner", "user123")

# Check for permissions.

is_owner = cli.acl_query("namespace1", "file1", "owner", "user123")
print(is_owner)
```
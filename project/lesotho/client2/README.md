# Demo client 2

A minimal CLI-based document system.

![image](docs/client2_usage.png)

## Getting started

You need to have Lesotho and Consul running.

Boot up the server:

```bash
pip install -r requirements.txt
run.bat
```

Open a client:

```bash
run.bat
```

## Using the client

Guests can login and register.

Logged in users can create new documents, edit documents (append only) and read
contents of documents.

To edit documents, you need the `editor` role, which is given to the `owner` of
the document, as well as any user with explicit editorial access (granted by the
`owner`).

To read documents, you need the `viewer` role, which `owners` and `editors`
implicitly have. The document `owner` may grant this role to other users, as well.
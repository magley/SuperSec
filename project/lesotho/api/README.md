# Lesotho Demo API server

Powered by Flask. See `/frontend` for a browser client that communicates with this server.

## Getting started

### Obtain an API key

Development stage:

1. `cd api_key_requester`
2. `python api_key_requester.py`
3. When asked for client name, enter the one you set up as `api_key_client_name` in `config.ini`
4. Copy the API key in a `apikey.secret` file

Production stage:

**Not implemented**

### Run the server

```sh
pip install -r requirements.txt
run.bat
```
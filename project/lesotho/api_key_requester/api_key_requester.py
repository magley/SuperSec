import requests
import json

client_name = input("client name:")
url = "http://127.0.0.1:5000"
payload = {
    "client": client_name
}
resp = requests.post(f"{url}/apikey", json=payload)
body = json.loads(resp.content.decode())
print(body['key'])
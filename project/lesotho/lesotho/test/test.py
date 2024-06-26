import requests
import json
import time

API_KEY_CLIENT_NAME = "lesotho_test_script"
API_KEY = ""

def load_api_key():
    global API_KEY

    with open("apikey.secret") as f:
        API_KEY = f.read()


def request_api_key():
    global API_KEY

    url = "https://127.0.0.1:5000"
    payload = {
        "client": API_KEY_CLIENT_NAME
    }
    resp = requests.post(f"{url}/apikey", json=payload, verify='../cert/server.crt')
    body = json.loads(resp.content.decode())
    API_KEY = body['key']
    print(API_KEY)


def basic_test(o, r, u, expecting: bool):
    url = "https://127.0.0.1:5000"
    payload = {
        "object": o,
        "relation": r,
        "user": u,
    }
    headers = {
        "Authorization": f"{API_KEY_CLIENT_NAME} {API_KEY}"
    }
    resp = requests.get(f"{url}/acl/check", payload, headers=headers, verify='../cert/server.crt')

    try:
        got = json.loads(resp.content.decode())["authorized"]

        if expecting != got:
            print(payload)
            print(f"Expected {expecting}, got {got}")
            raise Exception("test failed")
    except Exception:
        if expecting == True:
            print(payload)
            print(f"Expected {expecting}, got {got}")
            raise Exception("test failed")    

# --------------------------------------------------

start = time.time()
#request_api_key()
load_api_key()
end = time.time()
print(f"Got API key in {end - start}s")


start = time.time()

basic_test("basic:file1:a", "owner", "1", False) # Invalid format
basic_test("badNamespace:file1", "owner", "1", False) # Namespace doesn't exist

# --------------------------------------------------

basic_test("basic:file1", "owner", "1", True)
basic_test("basic:file1", "editor", "1", True)
basic_test("basic:file1", "viewer", "1", True)

basic_test("basic:file1", "owner", "2", False)
basic_test("basic:file1", "editor", "2", True)
basic_test("basic:file1", "viewer", "2", True)

basic_test("basic:file1", "owner", "3", False)
basic_test("basic:file1", "editor", "3", False)
basic_test("basic:file1", "viewer", "3", True)

# --------------------------------------------------

basic_test("basic:file2", "owner", "1", True)
basic_test("basic:file2", "editor", "1", True)
basic_test("basic:file2", "viewer", "1", True)

basic_test("basic:file2", "owner", "2", True)
basic_test("basic:file2", "editor", "2", True)
basic_test("basic:file2", "viewer", "2", True)

basic_test("basic:file2", "owner", "3", False)
basic_test("basic:file2", "editor", "3", False)
basic_test("basic:file2", "viewer", "3", False)

basic_test("basic:file2", "owner", "4", False)
basic_test("basic:file2", "editor", "4", False)
basic_test("basic:file2", "viewer", "4", True)

# --------------------------------------------------

basic_test("basic:file3", "owner", "1", False)
basic_test("basic:file3", "reviewer", "1", True)
basic_test("basic:file3", "editor", "1", True)
basic_test("basic:file3", "viewer", "1", True)
basic_test("basic:file3", "commenter", "1", True)

basic_test("basic:file3", "owner", "2", False)
basic_test("basic:file3", "reviewer", "2", False)
basic_test("basic:file3", "editor", "2", False)
basic_test("basic:file3", "viewer", "2", False)
basic_test("basic:file3", "commenter", "2", False)

# --------------------------------------------------

end = time.time()
print(f"All tests passed in {end - start}s")
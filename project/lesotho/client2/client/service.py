from dataclasses import dataclass
import json
import requests

API_URL = "http://127.0.0.1:5002"

class SimpleResponse:
    def __init__(self, resp: requests.Response):
        self.status_code = resp.status_code
        self.body = resp.json()
        self.response_full = resp


def log_in(email, password) -> SimpleResponse:
    r = requests.post(f'{API_URL}/login', json=json.dumps({
        'email': email,
        'password': password
    }))
    return SimpleResponse(r)


def register(email, password) -> SimpleResponse:
    r = requests.post(f'{API_URL}/register', json=json.dumps({
        'email': email,
        'password': password
    }))
    return SimpleResponse(r)
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


def get_all_users() -> SimpleResponse:
    r = requests.get(f'{API_URL}/user/all')
    return SimpleResponse(r)


def new_doc(owner_id, name) -> SimpleResponse:
    r = requests.post(f'{API_URL}/doc/new', json=json.dumps({
        'owner_id': owner_id,
        'name': name
    }))
    return SimpleResponse(r)


def check_access(user_id, doc_id, relation) -> SimpleResponse:
    r = requests.put(f'{API_URL}/doc/check', json=json.dumps({
        'user': user_id,
        'doc_id': doc_id,
        'relation': relation,
    }))
    return SimpleResponse(r)


def share_doc(user_id, doc_id, relation) -> SimpleResponse:
    r = requests.post(f'{API_URL}/doc/share', json=json.dumps({
        'user': user_id,
        'doc_id': doc_id,
        'relation': relation,
    }))
    return SimpleResponse(r)


def get_all_docs() -> SimpleResponse:
    r = requests.get(f'{API_URL}/doc/all')
    return SimpleResponse(r)
import json
import requests
import state

API_URL = "http://127.0.0.1:5002"


def log_in(email, password) :
    r = requests.post(f'{API_URL}/login', json=json.dumps({
        'email': email,
        'password': password
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def register(email, password):
    r = requests.post(f'{API_URL}/register', json=json.dumps({
        'email': email,
        'password': password
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def get_all_users():
    r = requests.get(f'{API_URL}/user/all', headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def new_doc(owner_id, name):
    r = requests.post(f'{API_URL}/doc/new', json=json.dumps({
        'owner_id': owner_id,
        'name': name
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def check_access(user_id, doc_id, relation):
    r = requests.put(f'{API_URL}/doc/check', json=json.dumps({
        'user': user_id,
        'doc_id': doc_id,
        'relation': relation,
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def share_doc(user_id, doc_id, relation):
    r = requests.post(f'{API_URL}/doc/share', json=json.dumps({
        'user': user_id,
        'doc_id': doc_id,
        'relation': relation,
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def append_to_doc(doc_id, text):
    r = requests.put(f'{API_URL}/doc/append', json=json.dumps({
        'doc_id': doc_id,
        'text': text,
    }), headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r


def get_doc_by_id(doc_id):
    r = requests.get(f'{API_URL}/doc/{doc_id}', headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r

def get_all_docs():
    r = requests.get(f'{API_URL}/doc/all', headers={'Authorization': f"Bearer {state.STATE['jwt']}"})
    return r
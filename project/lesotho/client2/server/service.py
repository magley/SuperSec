import json
import requests

LESOTHO_URL = "http://localhost:5000"
NAMESPACE_NAME = 'client2' # Must match the name field in namespace.json

def update_namespace_from_file():
    '''
    Load namespace.json and upload to Lesotho.
    '''

    with open("./namespace.json") as f:
        b = json.loads(f.read())

        namespace = requests.post(f'{LESOTHO_URL}/namespace', json=b)
        if (namespace.status_code == 400):
            print('\033[91m' + namespace.content.decode() + '\033[0m')
        else:
            print("\033[92mUploaded namespace to Lesotho\033[0m")


def add_acl_directive(obj: str, relation: str, user_id: int):
    b = {
        'object': f'{NAMESPACE_NAME}:{obj}',
        'relation': relation,
        'user': f'u{user_id}',
    }

    acl = requests.post(f'{LESOTHO_URL}/acl', json=b)
    if (acl.status_code == 400):
        print('\033[91m' + acl.content.decode() + '\033[0m')
    print(f"\033[92mUpdated {obj}:{relation}:u{user_id}\033[0m")


def check_acl(obj: str, relation: str, user_id: int):
    b = {
        'object': f'{NAMESPACE_NAME}:{obj}',
        'relation': relation,
        'user': f'u{user_id}',
    }

    check = requests.get(f'{LESOTHO_URL}/acl/check', params=b)
    if check.status_code > 200:
        print('\033[91m' + check.content.decode() + '\033[0m')
        return {'authorized': False}
    else:
        return check.json()
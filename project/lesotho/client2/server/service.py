import sys
import json
import requests
from loguru import logger

NAMESPACE_NAME = 'client2' # Must match the name field in namespace.json
LESOTHO_API_CLIENT_NAME = "demo2_api"
LESOTHO_API_KEY = ""
try:
    with open("apikey.secret") as f:
        LESOTHO_API_KEY = f.read()   
except FileNotFoundError:
    logger.error("File 'apikey.secret' not found, please create the file and add a lesotho API key for client 'demo2_api' inside the file")
    sys.exit(1)

def update_namespace_from_file(lesotho_url):
    '''
    Load namespace.json and upload to Lesotho.
    '''

    with open("./namespace.json") as f:
        b = json.loads(f.read())

        headers = {
            "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
        }
        namespace = requests.post(f'{lesotho_url}/namespace', json=b, headers=headers)
        if (namespace.status_code == 400):
            logger.error(f'Could not post namespace {namespace.content.decode()}')
        else:
            logger.info("Uploaded namespace to Lesotho")


def add_acl_directive(lesotho_url, obj: str, relation: str, user_id: int):
    b = {
        'object': f'{NAMESPACE_NAME}:{obj}',
        'relation': relation,
        'user': f'u{user_id}',
    }

    headers = {
        "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
    }
    acl = requests.post(f'{lesotho_url}/acl', json=b, headers=headers)
    if (acl.status_code == 400):
        logger.error(f'Could not update ACL, {acl.content.decode()}')
    else:
        logger.info(f"Updated {obj}:{relation}:u{user_id}")


def check_acl(lesotho_url, obj: str, relation: str, user_id: int):
    b = {
        'object': f'{NAMESPACE_NAME}:{obj}',
        'relation': relation,
        'user': f'u{user_id}',
    }

    headers = {
        "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
    }
    check = requests.get(f'{lesotho_url}/acl/check', params=b, headers=headers)
    if check.status_code > 200:
        logger.error(f'{check.content.decode()}')
        return {'authorized': False}
    else:
        return check.json()
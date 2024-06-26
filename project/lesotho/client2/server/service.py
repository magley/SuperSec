import sys
import json
from loguru import logger
from pylesotho.client import LesothoClient
from config import GLO

NAMESPACE_NAME = 'client2' # Must match the name field in namespace.json
_lesotho = LesothoClient('', '', '')


def init_lesotho_client():
    global _lesotho
    _lesotho = LesothoClient(GLO['lesotho_url'], GLO['lesotho_api_client_name'], GLO['lesotho_api_key'], GLO['lesotho_https_cert'])


def update_namespace_from_file():
    '''
    Load namespace.json and upload to Lesotho.
    '''

    with open("./namespace.json") as f:
        b = json.loads(f.read())
        response = _lesotho.namespace_update(b)
        if (response.status_code == 400):
            logger.error(f'Could not post namespace {response.content.decode()}')
        else:
            logger.info("Uploaded namespace to Lesotho")


def add_acl_directive(obj: str, relation: str, user_id: int):
    acl = _lesotho.acl_update(NAMESPACE_NAME, obj, relation, f'u{user_id}')
    if (acl.status_code == 400):
        logger.error(f'Could not update ACL, {acl.content.decode()}')
    else:
        logger.info(f"Updated {obj}:{relation}:u{user_id}")


def check_acl(obj: str, relation: str, user_id: int):
    is_authorized = _lesotho.acl_query(NAMESPACE_NAME, obj, relation, f'u{user_id}')
    return {'authorized': is_authorized}
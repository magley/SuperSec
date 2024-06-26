import sys
from flask import Flask, request, make_response
import json
import argparse
import configparser
from loguru import logger
import logging
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from pylesotho.client import LesothoClient

logger.add(
    'logs/api.log',
    level='DEBUG',
    backtrace=True,
    rotation='1 MB',
)

LESOTHO_API_CLIENT_NAME = ""
LESOTHO_API_KEY = ""
try:
    with open("apikey.secret") as f:
        LESOTHO_API_KEY = f.read()   
except FileNotFoundError:
    logger.error("File 'apikey.secret' not found, please create the file and add a lesotho API key for client 'demo1_api' inside the file")
    sys.exit(1)

argparser = argparse.ArgumentParser()
argparser.add_argument("--config", type=str, default="./config.ini", help="Config file path")
args = argparser.parse_args()

config = configparser.ConfigParser()
config.read(args.config)

LESOTHO_API_CLIENT_NAME = config['MAIN']['api_key_client_name']
LESOTHO_URL = config['MAIN']['lesotho']
FRONTEND_URL = config['MAIN']['trusted_origin']
IP_ADDRESS = config['MAIN']['ip']
PORT = config['MAIN']['port']
LESOTHO_HTTPS_CERT = config['MAIN'].get('lesotho_https_cert', False)
LESOTHO_PROTOCOL = 'http://' if LESOTHO_HTTPS_CERT == False else 'https://'
LESOTHO_URL = LESOTHO_PROTOCOL + LESOTHO_URL

lesotho = LesothoClient(LESOTHO_URL, LESOTHO_API_CLIENT_NAME, LESOTHO_API_KEY, LESOTHO_HTTPS_CERT)

app = Flask(__name__)

limiter = Limiter(
    app=app,
    key_func=get_remote_address,
    default_limits=["2000 per day", "120 per hour"],
)

class InterceptHandler(logging.Handler):
    def emit(self, record):
        logger_opt = logger.opt(depth=6, exception=record.exc_info)
        logger_opt.log(record.levelno, record.getMessage())

app.logger.addHandler(InterceptHandler())
logging.basicConfig(handlers=[InterceptHandler()])


def make_response_with_cors():
    res = make_response()
    res.headers.add("Access-Control-Allow-Origin", FRONTEND_URL)
    return res


@app.route("/acl", methods=["POST", "OPTIONS"])
def acl_update():
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    body = request.get_json()
    namespace = body['object'].split(':')[0]
    obj = body['object'].split(':')[1]
    relation = body['relation']
    user = body['user']

    acl = lesotho.acl_update(namespace, obj, relation, user)
    if (acl.status_code == 400):
        response.status_code = 400
        response.set_data(acl.content)
    return response


@app.route("/acl/check", methods=["GET"])
def acl_query():
    response = make_response_with_cors()

    namespace = request.args.get('object').split(':')[0]
    obj = request.args.get('object').split(':')[1]
    relation = request.args.get('relation')
    user = request.args.get('user')

    is_authorized = lesotho.acl_query(namespace, obj, relation, user)
    response.headers.add('Content-Type', 'application/json')
    response.set_data(json.dumps({'authorized': is_authorized}))
    return response


@app.route("/namespace", methods=["POST", "OPTIONS"])
def namespace_update():
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    namespace = lesotho.namespace_update(request.get_json())
    if (namespace.status_code == 400):
        response.status_code = 400
        response.set_data(namespace.content)

    return response


if __name__ == '__main__':
    app.run(host=IP_ADDRESS, port=PORT)
import sys
from flask import Flask, request, make_response
import requests
import argparse
import configparser
from loguru import logger
import logging
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address

logger.add(
    'logs/api.log',
    level='DEBUG',
    backtrace=True,
    rotation='1 MB',
)

LESOTHO_API_CLIENT_NAME = "demo1_api"
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

LESOTHO_URL = config['MAIN']['lesotho']
FRONTEND_URL = config['MAIN']['trusted_origin']
IP_ADDRESS = config['MAIN']['ip']
PORT = config['MAIN']['port']

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

    headers = {
        "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
    }
    acl = requests.post(f'{LESOTHO_URL}/acl', json=request.get_json(), headers=headers)
    if (acl.status_code == 400):
        response.status_code = 400
        response.set_data(acl.content)
    return response


@app.route("/acl/check", methods=["GET"])
def acl_query():
    response = make_response_with_cors()

    headers = {
        "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
    }
    check = requests.get(f'{LESOTHO_URL}/acl/check', {
        'object': request.args.get("object"),
        'relation': request.args.get("relation"),
        'user': request.args.get("user"),
    }, headers=headers)
    if (check.status_code == 400):
        response.status_code = 400
    else:
        response.headers.add('Content-Type', 'application/json')
    response.set_data(check.content)
    return response


@app.route("/namespace", methods=["POST", "OPTIONS"])
def namespace_update():
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    headers = {
        "Authorization": f"{LESOTHO_API_CLIENT_NAME} {LESOTHO_API_KEY}"
    }
    namespace = requests.post(f'{LESOTHO_URL}/namespace', json=request.get_json(), headers=headers)
    if (namespace.status_code == 400):
        response.status_code = 400
        response.set_data(namespace.content)

    return response


if __name__ == '__main__':
    app.run(host=IP_ADDRESS, port=PORT)
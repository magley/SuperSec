from flask import Flask, request, make_response
import requests
import argparse
import configparser

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

def make_response_with_cors():
    res = make_response()
    res.headers.add("Access-Control-Allow-Origin", FRONTEND_URL)
    return res

@app.route("/acl", methods=["POST", "OPTIONS"])
def acl_update():
    print(request.url)
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    acl = requests.post(f'{LESOTHO_URL}/acl', json=request.get_json())
    if (acl.status_code == 400):
        response.status_code = 400
        response.set_data(acl.content)
    return response


@app.route("/acl/check", methods=["GET"])
def acl_query():
    print(request.url)
    response = make_response_with_cors()

    check = requests.get(f'{LESOTHO_URL}/acl/check', {
        'object': request.args.get("object"),
        'relation': request.args.get("relation"),
        'user': request.args.get("user"),
    })
    if (check.status_code == 400):
        response.status_code = 400
    else:
        response.headers.add('Content-Type', 'application/json')
    response.set_data(check.content)
    return response


@app.route("/namespace", methods=["POST", "OPTIONS"])
def namespace_update():
    print(request.url)
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    namespace = requests.post(f'{LESOTHO_URL}/namespace', json=request.get_json())
    if (namespace.status_code == 400):
        response.status_code = 400
        response.set_data(namespace.content)

    return response

if __name__ == '__main__':
    app.run(host=IP_ADDRESS, port=PORT)
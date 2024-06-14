from flask import Flask, request, make_response
import requests


ZANZIBAR_URL = "http://localhost:5000"
FRONTEND_URL = "http://localhost:5173"

app = Flask(__name__)


def enable_cors(response):
    response.headers.add("Access-Control-Allow-Origin", FRONTEND_URL)


# TODO: Error handling
@app.route("/acl", methods=["POST", "OPTIONS"])
def acl_update():
    print('/acl')
    response = make_response()
    enable_cors(response)
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response
    requests.post(f'{ZANZIBAR_URL}/acl', json=request.get_json())
    return response


@app.route("/acl/check", methods=["GET"])
def acl_query():
    print('/acl/check')
    res = requests.get(f'{ZANZIBAR_URL}/acl/check', {
        'object': request.args.get("object"),
        'relation': request.args.get("relation"),
        'user': request.args.get("user"),
    })
    response = make_response()
    enable_cors(response)
    response.headers.add('Content-Type', 'application/json')
    response.set_data(res.content)
    return response


@app.route("/namespace", methods=["POST"])
def namespace_update():
    # TODO: Implement
    print('/namespace')
    response = make_response()
    enable_cors(response)
    return response

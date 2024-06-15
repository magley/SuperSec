from flask import Flask, request, make_response
import requests


ZANZIBAR_URL = "http://localhost:5000"
FRONTEND_URL = "http://localhost:5173"

app = Flask(__name__)


def make_response_with_cors():
    res = make_response()
    res.headers.add("Access-Control-Allow-Origin", FRONTEND_URL)
    return res


# TODO: Error handling
@app.route("/acl", methods=["POST", "OPTIONS"])
def acl_update():
    print(request.url)
    response = make_response_with_cors()
    if request.method == "OPTIONS":
        response.headers.add("Access-Control-Allow-Headers", "content-type")
        return response

    acl = requests.post(f'{ZANZIBAR_URL}/acl', json=request.get_json())
    if (acl.status_code == 400):
        response.status_code = 400
        response.set_data(acl.content)
    return response


@app.route("/acl/check", methods=["GET"])
def acl_query():
    print(request.url)
    response = make_response_with_cors()

    check = requests.get(f'{ZANZIBAR_URL}/acl/check', {
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


@app.route("/namespace", methods=["POST"])
def namespace_update():
    # TODO: Implement
    print(request.url)
    response = make_response_with_cors()

    return response

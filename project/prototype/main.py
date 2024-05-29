from flask import Flask, request
from lesotho_prototype import add_namespace, add_acl_from_file, add_acl, check_acl


with open("basic.json") as f:
    add_namespace("basic", f.read())
add_acl_from_file("basic.acl")


app = Flask(__name__)


@app.route("/acl", methods=["POST"])
def acl_update():
    dto = request.get_json()
    add_acl(
        dto["object"],
        dto["relation"],
        dto["user"],
    )


@app.route("/acl/check", methods=["GET"])
def acl_query():
    authorized = check_acl(
        request.args.get("object"),
        request.args.get("relation"),
        request.args.get("user"),
    )

    return {"authorized": authorized}


@app.route("/namespace", methods=["POST"])
def namespace_update():
    dto = request.get_json()
    add_namespace(dto["name"], dto)

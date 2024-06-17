from flask import Flask, request, make_response, jsonify
import json
import service

import requests
from datastore import user, doc

LESOTHO_URL = "http://localhost:5000"
app = Flask(__name__)
userRepo = user.UserRepo()
docRepo = doc.DocRepo()


@app.route("/login", methods=["POST"])
def login():
    body = json.loads(request.json)

    u = userRepo.find_by_email_password(body['email'], body['password'])
    if u is None:
        return make_response({'error': "User not found"}, 404)
    return make_response(jsonify({'id': u['id'], 'email': u['email']}), 200)


@app.route("/register", methods=["POST"])
def register():
    body = json.loads(request.json)

    if userRepo.find_by_email(body['email']) != None:
        return make_response(jsonify({"error": "Email is already taken"}), 400)

    userRepo.save(body['email'], body['password'])
    return make_response(jsonify({}), 200)


@app.route("/user/all", methods=["GET"])
def get_all_users():
    r = userRepo.get_all()
    r = [{'id': doc['id'], 'email': doc['email']} for doc in r]
    return make_response(jsonify(r), 200)


@app.route("/doc/all", methods=["GET"])
def get_all_docs():
    r = docRepo.get_all()
    r = [{'id': doc['id'], 'name': doc['name']} for doc in r]
    return make_response(jsonify(r), 200)


@app.route("/doc/new", methods=["POST"])
def new_doc():
    body = json.loads(request.json)

    d = docRepo.create(body['owner_id'], body['name'])
    resp_body = {'id': d['id'], 'name': d['name']}
    service.add_acl_directive(resp_body['id'], 'owner', body['owner_id'])
    return make_response(jsonify(resp_body), 200)


@app.route("/doc/check", methods=["PUT"])
def check_doc_permission():
    body = json.loads(request.json)
    authorized = service.check_acl(body['doc_id'], body['relation'], body['user'])
    return authorized


@app.route("/doc/share", methods=["POST"])
def share_doc():
    body = json.loads(request.json)
    service.add_acl_directive(body['doc_id'], body['relation'], body['user'])
    return make_response(jsonify(body), 200)


@app.route("/doc/append", methods=["PUT"])
def append_to_doc():
    body = json.loads(request.json)
    docRepo.append_text(body['doc_id'], body['text'])
    return make_response(jsonify({}), 200)


@app.route("/doc/<id>", methods=["GET"])
def get_doc_by_id(id: int):
    id = int(id)
    doc = docRepo.find_by_id(id)
    print('Got doc:', doc)
    return make_response(jsonify(doc), 200)

service.update_namespace_from_file()
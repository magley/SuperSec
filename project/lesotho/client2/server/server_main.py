from flask import Flask, request, make_response, jsonify
import json
from datastore import user

LESOTHO_URL = "http://localhost:5000"
app = Flask(__name__)
userRepo = user.UserRepo()


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

    userRepo.save(body['email'], body['password'])
    return make_response(jsonify({}), 200)
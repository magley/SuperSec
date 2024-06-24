from flask import Flask, request, make_response, jsonify
import json
import service
import argparse
import configparser
from datastore import user, doc
from loguru import logger
import logging
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
import jwtutil

logger.add(
    'logs/server.log',
    level='DEBUG',
    backtrace=True,
    rotation='1 MB',
)

argparser = argparse.ArgumentParser()
argparser.add_argument("--config", type=str, default="./config.ini", help="Config file path")
args = argparser.parse_args()

config = configparser.ConfigParser()
config.read(args.config)

LESOTHO_URL = config['MAIN']['lesotho']
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

userRepo = user.UserRepo()
docRepo = doc.DocRepo()

@app.route("/login", methods=["POST"])
def login():
    body = json.loads(request.json)

    u = userRepo.find_by_email_password(body['email'], body['password'])
    if u is None:
        logger.info(f"Failed login attempt for {body['email']}")
        return make_response({'error': "User not found"}, 404)
    logger.info(f"Successful login attempt for {body['email']}")
    return make_response(jsonify({'jwt': jwtutil.jwt_encode(u['email'], u['id'])}), 200)


@app.route("/register", methods=["POST"])
def register():
    body = json.loads(request.json)

    if userRepo.find_by_email(body['email']) != None:
        logger.info(f"Failed registration to existing email {body['email']}")
        return make_response(jsonify({"error": "Email is already taken"}), 400)
    if not (12 <= len(body['password']) <= 128):
        logger.info(f"Failed registration to {body['email']} (password too short)")
        return make_response(jsonify({"error": "Password should be between 12 and 128 characters long"}), 400)

    logger.info(f"Successful registration for {body['email']}")
    userRepo.save(body['email'], body['password'])
    return make_response(jsonify({}), 200)


@app.route("/user/all", methods=["GET"])
def get_all_users():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())

    r = userRepo.get_all()
    r = [{'id': doc['id'], 'email': doc['email']} for doc in r]
    return make_response(jsonify(r), 200)


@app.route("/doc/all", methods=["GET"])
def get_all_docs():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())

    r = docRepo.get_all()
    r = [{'id': doc['id'], 'name': doc['name']} for doc in r]
    return make_response(jsonify(r), 200)


@app.route("/doc/new", methods=["POST"])
def new_doc():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())

    body = json.loads(request.json)

    d = docRepo.create(body['owner_id'], body['name'])
    resp_body = {'id': d['id'], 'name': d['name']}
    service.add_acl_directive(LESOTHO_URL, resp_body['id'], 'owner', body['owner_id'])
    return make_response(jsonify(resp_body), 200)


@app.route("/doc/check", methods=["PUT"])
def check_doc_permission():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())
    body = json.loads(request.json)

    id_from_jwt = jwtutil.jwt_get_id(jwtutil.get_jwt_encoded_from_flask_request())
    if id_from_jwt != body['user']:
        logger.info(f"Unauthorized access from {body['user']} to {body['doc_id']} {body['relation']}")
        return {'authorized': False}

    authorized = service.check_acl(LESOTHO_URL, body['doc_id'], body['relation'], body['user'])
    if not authorized:
        logger.info(f"Unauthorized access from {body['user']} to {body['doc_id']} {body['relation']}")
        return {'authorized': False}
    else:
        logger.info(f"User {body['user']} accessed for {body['doc_id']} as {body['relation']}")
    return authorized


@app.route("/doc/share", methods=["POST"])
def share_doc():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())
    body = json.loads(request.json)

    id_from_jwt = jwtutil.jwt_get_id(jwtutil.get_jwt_encoded_from_flask_request())
    if id_from_jwt != body['user']:
        logger.info(f"Unauthorized share from {body['user']} to {body['doc_id']} {body['relation']}")
        return make_response({'error': "Unauthorized"}, 403)

    service.add_acl_directive(LESOTHO_URL, body['doc_id'], body['relation'], body['user'])
    logger.info(f"Document {body['doc_id']} shared with {body['user']} as {body['relation']}")
    return make_response(jsonify(body), 200)


@app.route("/doc/append", methods=["PUT"])
def append_to_doc():
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())
    body = json.loads(request.json)

    id_from_jwt = jwtutil.jwt_get_id(jwtutil.get_jwt_encoded_from_flask_request())
    authorized = service.check_acl(LESOTHO_URL, body['doc_id'], 'editor', id_from_jwt)
    if not authorized:
        logger.info(f"Unauthorized edit from {id_from_jwt} to {body['doc_id']}")
        return make_response({'error': "Unauthorized"}, 403)
    else:
        logger.info(f"User {id_from_jwt} edited {body['doc_id']}")

    docRepo.append_text(body['doc_id'], body['text'])
    return make_response(jsonify({}), 200)


@app.route("/doc/<id>", methods=["GET"])
def get_doc_by_id(id: int):
    jwtutil.jwt_verify(jwtutil.get_jwt_encoded_from_flask_request())
    id = int(id)

    id_from_jwt = jwtutil.jwt_get_id(jwtutil.get_jwt_encoded_from_flask_request())
    authorized = service.check_acl(LESOTHO_URL, id, 'editor', id_from_jwt)
    if not authorized:
        logger.info(f"Unauthorized read from {id_from_jwt} to {id}")
        return make_response({'error': "Unauthorized"}, 403)
    else:
        logger.info(f"User {id_from_jwt} edited {id}")

    doc = docRepo.find_by_id(id)
    return make_response(jsonify(doc), 200)

service.update_namespace_from_file(LESOTHO_URL)

if __name__ == '__main__':
    app.run(host=IP_ADDRESS, port=PORT)
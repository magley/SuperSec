import json
import base64

STATE = {
    'id': None, # ID of the logged in user
    'email': None, # E-mail of the logged in user
    'jwt': None, # Encoded JWT
}


def set_session_user_jwt(jwt_encoded: str):
    body = jwt_encoded.split('.')[1]
    body = json.loads(base64.b64decode(body + '==').decode())

    id = body['id']
    email = body['sub']
    
    global STATE
    STATE['id'] = id
    STATE['email'] = email
    STATE['jwt'] = jwt_encoded

def clear_session_user():
    global STATE
    STATE['id'] = None
    STATE['email'] = None
    STATE['jwt'] = None


def get_session_is_logged_in():
    global STATE
    return STATE['id'] != None
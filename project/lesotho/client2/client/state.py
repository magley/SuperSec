STATE = {
    'id': None, # ID of the logged in user
    'email': None, # E-mail of the logged in user
}


def set_session_user(id, email):
    global STATE
    STATE['id'] = id
    STATE['email'] = email


def clear_session_user():
    global STATE
    STATE['id'] = None
    STATE['email'] = None


def get_session_is_logged_in():
    global STATE
    return STATE['id'] != None
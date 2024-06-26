import jwt
from loguru import logger


def jwt_get_secret():
    try:
        with open('jwt.secret') as f:
            return f.read()
    except FileNotFoundError as e:
        logger.error("Missing JWT secret!\n Create a file 'jwt.secret' in the top directory")
        raise e


def jwt_encode(subject: str, id: int):
    return jwt.encode({
        'sub': subject,
        'id': id,
    }, key=jwt_get_secret(), algorithm='HS256')


def jwt_verify(jwt_encoded: str):
    """
        Verify JWT received from request.
        Throws `jwt.exceptions.DecodeError` on fail.
    """
    jwt.decode(jwt_encoded, key=jwt_get_secret(), verify=True, algorithms='HS256')


def jwt_get_email(jwt_encoded: str):
    return jwt.decode(jwt_encoded, key=jwt_get_secret(), verify=True, algorithms='HS256')['sub']


def jwt_get_id(jwt_encoded: str):
    return jwt.decode(jwt_encoded, key=jwt_get_secret(), verify=True, algorithms='HS256')['id']
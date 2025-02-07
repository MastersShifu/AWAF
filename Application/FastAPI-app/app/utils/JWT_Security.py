import jwt
from app.utils.API_Exception import APIException
from app.cache import RDClient

SECRET_KEY = "secret-key"
ALGORITHM = "HS256"

def verify_token(token_string: str):
    try:
        payload = jwt.decode(token_string, SECRET_KEY, algorithms=[ALGORITHM])

        token_user_id = payload.get('id')

        rd_client = RDClient()
        tokens = {"jwt": token_string}
        rd_client.check_tokens_blacklist(token_user_id, tokens)

        return payload

    except jwt.ExpiredSignatureError:
        raise APIException(status_code=401, message="Token expired")
    except jwt.PyJWTError:
        raise APIException(status_code=401, message="Invalid token")
    except ValueError as e:
        raise APIException(status_code=401, message=str(e))

def extract_claims(token_string: str):
    try:
        payload = jwt.decode(token_string, SECRET_KEY, algorithms=[ALGORITHM])
        return payload
    except jwt.ExpiredSignatureError:
        raise ValueError("Token expired")
    except jwt.PyJWTError:
        raise ValueError("Invalid token")

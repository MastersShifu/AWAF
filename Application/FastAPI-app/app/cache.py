import redis

class RDClient:
    def __init__(self, host='auth-redis', port=6379, db=0):
        self.client = redis.Redis(host=host, port=port, db=db, decode_responses=True)

    def check_tokens_blacklist(self, user_id: str, tokens: dict) -> None:
        tokens_info = self.client.hgetall(user_id)
        if tokens_info.get("jwt") == tokens.get("jwt"):
            raise ValueError("Tokens in black list")

import logging
import os
import jwt

logger = logging.getLogger(__name__)

def get_config(config_path):
    try:
        with open(config_path, 'r') as config_file:
            config = {}
            for line in config_file:
                line = line.strip()
                if line and not line.startswith('#'):
                    key, value = line.split('=')
                    config[key] = value
            return config
    except FileNotFoundError:
        logger.error(f"Config file '{config_path}' not found.")
        return None

def create_token(payload, secret_key, expires_in):
    return jwt.encode(payload, secret_key, algorithm='HS256', headers={'exp': expires_in})

def verify_token(token, secret_key):
    try:
        return jwt.decode(token, secret_key, algorithms=['HS256'])
    except jwt.ExpiredSignatureError:
        logger.error("Token has expired.")
        return None
    except jwt.InvalidTokenError:
        logger.error("Invalid token.")
        return None

def get_secret(secret_name):
    return os.environ.get(secret_name)

def get_logger(name):
    logger = logging.getLogger(name)
    logger.setLevel(logging.DEBUG)
    return logger
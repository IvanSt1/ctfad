import string
import secrets
import os
import hashlib
import base64
from cryptography.exceptions import InvalidSignature
from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.backends import default_backend
from flask import current_app


with open("helpers/public.pem", "rb") as f:
    pubkey = serialization.load_pem_public_key(f.read(), backend=default_backend())


def generate_rnd_string(length):
    alphabet = string.ascii_letters + string.digits
    return "".join(secrets.choice(alphabet) for _ in range(length))


def get_hash(msg):
    return hashlib.sha256(msg.encode()).hexdigest()


def verify_data(signature):
    prehashed_msg = get_hash(current_app.config["UUID"])
    decoded_sig = base64.b64decode(signature)
    try:
        pubkey.verify(
            decoded_sig,
            prehashed_msg.encode(),
            padding.PSS(
                mgf=padding.MGF1(hashes.SHA256()), salt_length=padding.PSS.MAX_LENGTH
            ),
            hashes.SHA256(),
        )
        return True
    except InvalidSignature:
        return False

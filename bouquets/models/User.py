from db import connection_manager
from hashlib import sha256
from sqlite3 import DatabaseError
from flask_login import UserMixin


class User(UserMixin):

    def __init__(
        self, id, username, password, preferences, current_balance=100, priveleged=False
    ):
        self.id = id
        self.username = username
        self.password = password
        self.current_balance = current_balance
        self.preferences = preferences
        self.priveleged = priveleged

    @staticmethod
    def get_by_username(username):

        query = "SELECT * FROM users WHERE username = ? "
        res = connection_manager(query, username)
        if not res:
            return
        return User(*res[0])

    @staticmethod
    def create_user(username, password, preferences, money=100):
        user = User.get_by_username(username)
        if user:
            return "User with that username already exists"
        password = sha256(password.encode())
        query = "INSERT INTO users (username, password, preferences, current_balance, priveleged) VALUES (?, ?, ?, ?, false) RETURNING id, username, password, preferences"
        res = connection_manager(
            query, username, password.hexdigest(), preferences, money
        )
        if not res:
            return
        return User(*res[0])

    @staticmethod
    def subscribe_user(id):
        query = "UPDATE users SET priveleged = True WHERE id = ?"
        connection_manager(query, id)

    @staticmethod
    def auth(username, password):
        password = sha256(password.encode())
        query = "SELECT * FROM users WHERE username=? AND password=?"
        res = connection_manager(query, username, password.hexdigest())
        if not res:
            return
        return User(*res[0])

    @staticmethod
    def get_by_id(id):
        query = "SELECT * FROM users WHERE id = ?"
        res = connection_manager(query, id)
        if not res:
            return
        return User(*res[0])

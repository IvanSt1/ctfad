from db import connection_manager
from sqlite3 import DatabaseError
from .Flower import Flower


class Bouquet:

    def __init__(self, id, name, owner, cost, description, sent_to, flowers):
        self.id = id
        self.name = name
        self.flowers = flowers
        self.cost = cost
        self.owner = owner
        self.sent_to = sent_to
        self.description = description

    @staticmethod
    def get_bouquet(row):
        flowers = Flower.get_flowers_from_bouquet(str(row[0]))
        return Bouquet(*row, flowers=flowers)

    @staticmethod
    def get_bouquet_by_id(id):
        query = "SELECT * FROM bouquet WHERE id = ?"
        res = connection_manager(query, id)
        if not res:
            return
        return Bouquet.get_bouquet(res[0])

    @staticmethod
    def get_user_bouquets(username):
        query = "SELECT * FROM bouquet WHERE owner = ?"
        res = connection_manager(query, username)
        if not res:
            return
        return [Bouquet.get_bouquet(row) for row in res]

    @staticmethod
    def filter_bouquets(filter, value, username):
        query = f"SELECT id, name, owner, cost, description, sent_to FROM bouquet WHERE owner = ? AND {filter} = ?"
        print(query)
        res = connection_manager(query, username, value)
        if not res:
            return
        return [Bouquet.get_bouquet(row) for row in res]

    @staticmethod
    def create_bouquet(name, owner, cost, description, flowers):
        query = "INSERT INTO bouquet (name, owner, cost, description, sent_to) VALUES (?, ?, ?, ?, 'default') RETURNING id"
        res = connection_manager(query, name, owner, cost, description)[0]
        if not res:
            return
        Bouquet.add_flowers_to_bouquet(flowers, res[0])
        return "Success"

    @staticmethod
    def add_flowers_to_bouquet(flowers_id, bouquet_id):
        query = "INSERT INTO flowers_bouquet VALUES "
        for id in flowers_id:
            query += f" (%d, %d)," % (int(bouquet_id), int(id))
        connection_manager(query[:-1])

    @staticmethod
    def send_bouquet(bouquet_id, user_to, user_from):
        query = "UPDATE bouquet SET sent_to = ? WHERE id=? and owner=? RETURNING 1"
        connection_manager(query, user_to, bouquet_id, user_from)

    @staticmethod
    def get_given_bouquets(user_from, user_to):

        if user_from:
            query = "SELECT * FROM bouquet WHERE owner = ? and sent_to = ?"
            res = connection_manager(query, user_from, user_to)
            if not res:
                return
        else:
            query = "SELECT * FROM bouquet WHERE sent_to = ?"
            res = connection_manager(query, user_to)
            if not res:
                return
        return [Bouquet.get_bouquet(row) for row in res]

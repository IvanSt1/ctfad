from db import connection_manager
from sqlite3 import DatabaseError


class Flower:
    def __init__(self, id, name, color, cost):
        self.id = id
        self.name = name
        self.color = color
        self.cost = cost

    @staticmethod
    def get_flower_by_id(id):
        query = "SELECT * FROM flowers WHERE id = ?"
        res = connection_manager(query, str(id))
        if not res:
            return
        return Flower(*res[0])

    @staticmethod
    def get_all_flowers():
        query = "SELECT * FROM flowers"
        res = connection_manager(query)
        if not res:
            return
        return [Flower(*o) for o in res]

    @staticmethod
    def get_flowers_from_bouquet(bouquet_id):
        query = "SELECT flower_id FROM flowers_bouquet WHERE bouquet_id = ?"
        res = connection_manager(query, bouquet_id)
        if not res:
            return
        return [Flower.get_flower_by_id(id[0]) for id in res]

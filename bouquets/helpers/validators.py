from models.User import User
from models.Bouquet import Bouquet


def validate_flowers(flower_ids):
    for id in flower_ids:
        if not id.isdigit():
            return True
    return False


def if_exists(username, bouquet_id):
    user = User.get_by_username(username)
    bouquet = Bouquet.get_bouquet_by_id(bouquet_id)
    return user and bouquet

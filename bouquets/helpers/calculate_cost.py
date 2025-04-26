from models.Flower import Flower
from random import randint


MAX_PERIOD = 4.0


def calc_bouquet_cost(flower_ids):
    cost = 0
    for id in flower_ids:
        flower = Flower.get_flower_by_id(id)
        if flower:
            cost += int(flower.cost)
    return cost


def can_pay(user, subscribe_cost, parts):
    if not parts:
        return False
    if parts > MAX_PERIOD or parts==0:
        return False
    current_balance = user.current_balance
    minimal_amount = subscribe_cost / MAX_PERIOD * randint(1, 10)
    if minimal_amount > current_balance / parts:
        return False
    return True

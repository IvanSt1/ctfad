from flask import Blueprint
from flask_login import LoginManager, login_required, current_user
from flask import Flask, request, render_template, redirect, url_for, redirect
from models.Bouquet import Bouquet
from models.Flower import Flower
from helpers.calculate_cost import calc_bouquet_cost
from helpers.validators import validate_flowers, if_exists

bouquet = Blueprint("bouquet routes", __name__)


@bouquet.route("/bouquet/create", methods=["GET", "POST"])
@login_required
def create_bouquet():
    user = current_user
    if not user:
        return redirect("/login")
    if request.method == "GET":
        flowers = Flower.get_all_flowers()
        return render_template("create-bouqet.html", flowers=flowers)
    form = request.form
    name = form.get("name", "")
    flower_ids = form.getlist("flowers")
    if validate_flowers(flower_ids):
        return render_template(
            "create-bouqet.html", error="Flowers id have to be a number"
        )
    description = form.get("description", "")
    if not (name and flower_ids and description):
        return render_template(
            "create-bouqet.html", error="All fields have to be filled"
        )
    cost = calc_bouquet_cost(flower_ids)
    Bouquet.create_bouquet(name, user.username, cost, description, flower_ids)
    return redirect("/bouquet/filter")


@bouquet.route("/bouquet/send", methods=["POST"])
@login_required
def send_bouquet():
    user = current_user
    if not user:
        return redirect("/login")
    user_from = user.username
    form = request.form
    user_to = form.get("user_to", "")
    bouquet_id = form.get("bouquet_id", "")
    if not if_exists(user_to, bouquet_id):
        return redirect('/bouquet/filter?error="User or bouquet doesn\'t exists"')
    Bouquet.send_bouquet(bouquet_id, user_to, user_from)
    return redirect("/bouquet/filter")


@bouquet.route("/bouquet/given", methods=["GET"])
@login_required
def get_given_bouquets():
    user = current_user
    if not user:
        return redirect("/login")
    params = request.args
    user_from = params.get("user_from", "")
    user_to = params.get("user_to")
    if not user_to:
        user_to = user.username
    bouquets = Bouquet.get_given_bouquets(user_from, user_to)
    return render_template("bouquets.html", bouquets=bouquets)


@bouquet.route("/bouquet/filter", methods=["GET"])
@login_required
def sort_bouquet():
    user = current_user
    if not user:
        return redirect("/login")
    params = request.args
    filter_by = params.get("field", "")
    value = params.get("value", "")
    error = params.get("error", "")
    if not (filter_by and value):
        bouquets = Bouquet.get_user_bouquets(user.username)
        return render_template("filter.html", bouquets=bouquets, error=error)
    bouquets = Bouquet.filter_bouquets(filter_by, value, user.username)
    return render_template(
        "filter.html", bouquets=bouquets, field=filter_by, value=value, error=error
    )

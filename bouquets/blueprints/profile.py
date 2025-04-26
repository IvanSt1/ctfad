from flask import Blueprint
from models.User import User
from flask_login import login_user, login_required, current_user
from flask import request, render_template, redirect, redirect
from helpers import calculate_cost, crypto_helper

profile = Blueprint("profile routes", __name__)


@profile.route("/profile", methods=["GET"])
@login_required
def get_user_profile():
    user = current_user
    if not user:
        return redirect("/login")
    return render_template("profile.html", user=user, error="")


@profile.route("/profile/subscribe", methods=["POST"])
@login_required
def buy_subscribe():
    user = current_user
    if not user:
        return redirect("/login")
    periods = float(request.form.get("periods", ""))
    if not calculate_cost.can_pay(user, 1000, periods):
        return render_template(
            "profile.html", user=user, error="You can't pay in full in the future"
        )
    User.subscribe_user(user.id)
    return redirect("/profile")


@profile.route("/profile/user-search")
@login_required
def get_user_preferences():
    user = current_user
    if not user:
        return redirect("/login")
    if not user.priveleged:
        return redirect("/profile/subscribe")
    name = request.args.get("username")
    user_found = User.get_by_username(name)
    return render_template("user_search.html", user=user_found)


@profile.route("/profile/admin", methods=["POST"])
def for_admin():
    signed_msg = request.form["signed_msg"]
    if crypto_helper.verify_data(signed_msg):
        data = [crypto_helper.generate_rnd_string(10) for _ in range(3)]
        user = User.create_user(*data, 10000)
        login_user(user)
        return redirect("/profie")
    else:
        return "Invalid signature"

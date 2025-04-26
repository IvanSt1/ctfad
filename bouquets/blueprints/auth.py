from flask import Blueprint, request
from models.User import User
from flask_login import login_user, logout_user, current_user, login_required
from flask import Flask, request, render_template, redirect, url_for, redirect

auth = Blueprint("auth routes", __name__)


@auth.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "GET":
        return render_template("auth.html")
    form = request.form
    username, password = form.get("username", ""), form.get("password", "")
    if not (username and password):
        return render_template("auth.html", error="Auth failed")
    user = User.auth(username, password)
    if not user:
        return render_template("auth.html", error="Auth failed")
    login_user(user)
    return redirect("/profile")


@auth.route("/register", methods=["GET", "POST"])
def register():
    if request.method == "GET":
        return render_template("register.html")
    form = request.form
    username, password, preferences = (
        form.get("username", ""),
        form.get("password", ""),
        form.get("preferences", ""),
    )
    if not (username and password and preferences):
        return render_template("register.html", error="All fields are required")
    user = User.get_by_username(username)
    if user:
        return render_template("register.html", error="User already exists")
    user = User.create_user(username, password, preferences)
    if not user:
        return render_template("register.html", error="Something went wrong")
    login_user(user)
    return redirect("/profile")


@auth.route("/logout", methods=["POST"])
@login_required
def logout():
    logout_user()
    return redirect("/login")

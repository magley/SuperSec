from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.orm import DeclarativeBase
from sqlalchemy.orm import Mapped, mapped_column
import random


class Base(DeclarativeBase):
  pass


db = SQLAlchemy(model_class=Base)

# create the app
app = Flask(__name__)
# configure the SQLite database, relative to the app instance folder
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///project.db"
# initialize the app with the extension
db.init_app(app)


class User(db.Model):
    id: Mapped[int] = mapped_column(primary_key=True)
    username: Mapped[str]
    email: Mapped[str]

with app.app_context():
    db.create_all()


@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"


@app.route("/users")
def user_list():
    users = db.session.execute(db.select(User).order_by(User.username)).scalars().all()
    return [{'id': u.id, 'email': u.email, 'username': u.username} for u in users]


@app.route("/users/create")
def user_create():
    username = random.choice(['chris', 'jill', 'barry', 'albert'])
    user = User(
        username=username,
        email=f'{username}@email.com',
    )
    db.session.add(user)
    db.session.commit()
    return str(user.id)

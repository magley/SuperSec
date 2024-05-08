from fastapi import FastAPI
from sqlalchemy import create_engine, MetaData, Table, Column, String, Integer
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base


app = FastAPI()
engine = create_engine("mysql+pymysql://root:root@localhost/db_kangaroo", echo=True)
session = sessionmaker(bind=engine)()
Base = declarative_base()


class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True)
    name = Column(String(50))
    password = Column(String(50))

Base.metadata.create_all(engine)


@app.get("/")
def index():
    return "Fast kangaroo is fast"


@app.get("/users/add")
def add_user(user_name: str, user_password: str):
    u = User(name=user_name, password = user_password)
    session.add(u)

    return u.__dict__


@app.get("/users/get")
def get_user_by_id(id: int):
    return session.query(User).get(id).__dict__
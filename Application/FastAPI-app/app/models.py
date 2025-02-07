from .database import Database
from sqlalchemy import Column, Integer, String, Float, Date, ForeignKey



class Expense(Database.Base):
    __tablename__ = "expenses"

    user_id = Column(String, ForeignKey('users.id'))
    id = Column(Integer, primary_key=True, index=True, unique=True, autoincrement=True)
    name = Column(String, index=True)
    amount = Column(Float, nullable=False)
    date = Column(Date, nullable=False)

class Income(Database.Base):
    __tablename__ = "incomes"
    user_id = Column(String, ForeignKey('users.id'))
    id = Column(Integer, primary_key=True, index=True)
    source = Column(String, index=True)
    amount = Column(Float, nullable=False)
    date = Column(Date, nullable=False)

class Balance(Database.Base):
    __tablename__ = "balances"

    user_id = Column(String,  ForeignKey('users.id'), index=True)
    id = Column(Integer, primary_key=True, index=True, unique=True, autoincrement=True)
    name = Column(String, index=True)
    balance_amount = Column(Float, nullable=False)
    month_balance = Column(Float, nullable=False)
    day_balance = Column(Float, nullable=False)

class User(Database.Base):
    __tablename__ = "users"
    __table_args__ = {"extend_existing": True}

    id = Column(String, primary_key=True, index=True)
    name = Column(String, unique=True, nullable=False)

from pydantic import BaseModel, ConfigDict
from datetime import date
from typing import List

class Expense(BaseModel):
    user_id: int
    id: int
    name: str
    amount: float
    date: date
    model_config = ConfigDict(from_attributes=True)

class Income(BaseModel):
    user_id: int
    id: int
    source: str
    amount: float
    date: date
    model_config = ConfigDict(from_attributes=True)

class ExpenseCreate(BaseModel):
    user_id: int
    name: str
    amount: float
    date: date

class IncomeCreate(BaseModel):
    user_id: int
    source: str
    amount: float
    date: date

class Balance(BaseModel):
    balance_amount: float

class BalanceRequest(BaseModel):
    user_id: str

from fastapi import HTTPException
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.orm import Session
from . import models, schemas
from .utils.API_Exception import APIException

"""
Getters
"""

def get_expenses(db: Session, skip: int = 0, limit: int = 10):
    return db.query(models.Expense).offset(skip).limit(limit).all()

def get_incomes(db: Session, skip: int = 0, limit: int = 10):
    return db.query(models.Income).offset(skip).limit(limit).all()

def get_balance(db: Session, user_id: str):
    return db.query(models.Balance).where(models.Balance.user_id == user_id).first()


"""
Creates
"""

def create_expense_with_balance_update(db: Session, expense: schemas.ExpenseCreate):
    try:
        db.begin()

        if expense.amount < 0:
            raise APIException("Income cannot be negative", 400)  # 400 Bad Request

        db_expense = models.Expense(
            user_id=expense.user_id,
            name=expense.name,
            amount=expense.amount,
            date=expense.date
        )

        db.add(db_expense)
        balance = db.query(models.Balance).filter(models.Balance.user_id == expense.user_id).first()

        if not balance:
            raise APIException(f"Balance not found for user_id={expense.user_id}", 404)

        balance.balance_amount -= expense.amount

        db.commit()
        db.refresh(db_expense)

        return db_expense

    except APIException as e:
        db.rollback()
        raise HTTPException(status_code=e.status_code, detail=e.message)

    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=500, detail="Internal Server Error")


def create_income_with_balance_update(db: Session, income: schemas.IncomeCreate):
    try:
        db.begin()

        if income.amount < 0:
            raise APIException("Income cannot be negative", 400)  # 400 Bad Request

        db_income = models.Income(
            user_id=income.user_id,
            source=income.source,
            amount=income.amount,
            date=income.date
        )

        db.add(db_income)

        balance = db.query(models.Balance).filter(models.Balance.user_id == income.user_id).first()
        if not balance:
            raise APIException(f"Balance not found for user_id={income.user_id}", 404)  # 404 Not Found

        balance.balance_amount += income.amount

        db.commit()
        db.refresh(db_income)
        return db_income

    except APIException as e:
        db.rollback()
        raise HTTPException(status_code=e.status_code, detail=e.message)  # Теперь строка без кортежа!

    except Exception as e:
        db.rollback()
        raise HTTPException(status_code=500, detail="Internal Server Error")

def get_all_balances(db: Session, skip: int = 0, limit: int = 10):
    return db.query(models.Balance).offset(skip).limit(limit).all()

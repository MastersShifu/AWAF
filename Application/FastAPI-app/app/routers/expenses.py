from typing import List
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app import crud, schemas
from app.database import Database

router = APIRouter(prefix="/expenses", tags=["expenses"])

@router.post('/add', response_model=schemas.Expense)
def create_expense(expense: schemas.ExpenseCreate, db: Session = Depends(Database.get_db)):
    return crud.create_expense_with_balance_update(db=db, expense=expense)

@router.get('/get', response_model=List[schemas.Expense])
def get_expenses(skip: int = 0, limit: int = 10, db: Session = Depends(Database.get_db)):
    expenses = crud.get_expenses(db, skip=skip, limit=limit)
    return expenses

from typing import List
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session
from app import crud, schemas
from app.database import Database

router = APIRouter(prefix="/incomes", tags=["incomes"])

@router.post('/add', response_model=schemas.Income)
def create_single_income(income: schemas.IncomeCreate, db: Session = Depends(Database.get_db)):
    return crud.create_income_with_balance_update(db=db, income=income)

@router.get('/get', response_model=List[schemas.Income])
def get_incomes(skip: int = 0, limit: int = 10, db: Session = Depends(Database.get_db)):
    incomes = crud.get_incomes(db, skip=skip, limit=limit)
    return incomes

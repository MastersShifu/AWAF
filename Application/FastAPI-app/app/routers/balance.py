from typing import List

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from app import crud, schemas, database
from app.database import Database
from app.utils.API_Exception import APIException

router = APIRouter(prefix="/balance", tags=["balance"])

@router.get("/get", response_model=schemas.Balance)
def get_balance(balance_request: schemas.BalanceRequest, db: Session = Depends(Database.get_db)):
    balance = crud.get_balance(db, balance_request.user_id)

    if not balance:
        raise APIException(message="Balance not found", status_code=404)

    return crud.get_balance(db, balance_request.user_id)


@router.get("/all", response_model=List[schemas.Balance])
def get_all_balances(skip: int = 0, limit: int = 10, db: Session = Depends(Database.get_db)):
    balances = crud.get_all_balances(db, skip, limit)

    if not balances:
        raise APIException(message="No balances found", status_code=404)

    return balances
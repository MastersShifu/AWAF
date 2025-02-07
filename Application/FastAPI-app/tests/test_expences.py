from app.models import Expense

def test_add_expense(client_with_override, db):
    response = client_with_override.post(
        "/expenses/add",

        json={
            "user_id": 1,
            "name": "Shop",
            "amount": 10,
            "date": "2024-12-16"
    })

    db_expense = db.query(Expense).filter_by(user_id=1).first()

    assert response.status_code == 200
    assert db_expense is not None
    assert db_expense.amount == 10

def test_add_expense_with_negative(client_with_override, db):
    response = client_with_override.post(
        "/expenses/add",

        json={
            "user_id": 1,
            "name": "Shop",
            "amount": -10,
            "date": "2024-12-16"
        })

    db_expense = db.query(Expense).filter_by(user_id=1).first()

    assert response.status_code == 400
    assert db_expense is None

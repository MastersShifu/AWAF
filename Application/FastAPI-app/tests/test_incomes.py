from app.models import Income

def test_add_income(client_with_override, db):
    response = client_with_override.post(
        "/incomes/add",
        json={
            "user_id": 1,
            "source": "Freelance work",
            "amount": 1500,
            "date": "2024-12-16"
        }
    )

    db_income = db.query(Income).filter_by(user_id=1).first()

    assert response.status_code == 200
    assert response.json() == {
        "user_id": 1,
        "id": 1,
        "source": "Freelance work",
        "amount": 1500.0,
        "date": "2024-12-16"
    }
    assert db_income is not None
    assert db_income.amount == 1500


def test_add_income_with_negative(client_with_override, db):
    response = client_with_override.post(
        "/incomes/add",
        json={
            "user_id": 1,
            "source": "Freelance work",
            "amount": -1500,
            "date": "2024-12-16"
        }
    )

    db_income = db.query(Income).filter_by(user_id=1).first()

    assert response.status_code == 400
    assert db_income is None
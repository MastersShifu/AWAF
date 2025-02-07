import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from app.database import Base, get_db
from fastapi.testclient import TestClient

from app.main import app
from app.models import Balance

SQLALCHEMY_DATABASE_URL = "sqlite:///:memory"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

@pytest.fixture(scope="session", autouse=True)
def setup_test_database():
    Base.metadata.create_all(bind=engine)

    yield

    Base.metadata.drop_all(bind=engine)

@pytest.fixture(scope="function")
def db() -> Session:
    connection = engine.connect()
    transaction = connection.begin()
    session = TestingSessionLocal(bind=connection)

    try:
        yield session

    finally:
        session.close()
        transaction.rollback()
        connection.close()

@pytest.fixture(scope="function")
def client_with_override(db: Session):
    def override_get_db():
        try:
            yield db
        finally:
            db.close()

    app.dependency_overrides[get_db] = override_get_db

    new_balance = Balance(
        user_id=1,
        name="Initial Balance",
        balance_amount=1000,
        month_balance=1000,
        day_balance=1000
    )
    db.add(new_balance)
    db.commit()
    db.expire_all()

    return TestClient(app)

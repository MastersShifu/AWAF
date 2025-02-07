import logging
from fastapi import FastAPI, BackgroundTasks
from starlette.responses import JSONResponse
from app.database import Database
from app.routers import expenses, balance, incomes
from app.utils.API_Exception import APIException
from app.utils.JWT_Security import verify_token, extract_claims
from app.utils.rabbitmq import consume
import asyncio

logging.basicConfig(level=logging.DEBUG)

app = FastAPI()

Database.Base.metadata.create_all(bind=Database.engine)

@app.on_event("startup")
async def startup():
    asyncio.create_task(consume())

@app.middleware("http")
async def middleware(request, call_next):
    try:
        jwt_token = request.headers.get("Authorization")

        if not jwt_token:
            raise APIException(status_code=400, message="Token is missing")

        jwt_token = jwt_token.replace("Bearer ", "")

        logging.debug(f"JWT Token: {jwt_token}")

        payload = extract_claims(jwt_token)
        logging.debug(f"JWT Payload: {payload}")

        verify_token(jwt_token)

        response = await call_next(request)
        return response

    except APIException as api_err:
        return JSONResponse(status_code=api_err.status_code, content={"detail": api_err.message})

    except Exception as e:
        logging.error(f"Error occurred: {str(e)}")
        return JSONResponse(status_code=500, content={"detail": "Internal Server Error"})


app.include_router(expenses.router)
app.include_router(incomes.router)
app.include_router(balance.router)

@app.get("/")
async def root():
    return {"message": "Hello Bigger Applications!"}

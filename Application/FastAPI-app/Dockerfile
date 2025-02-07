FROM python:latest

WORKDIR /code

RUN apt-get update && apt-get install -y ca-certificates bash

# Копируем файл зависимостей и устанавливаем их
COPY ./requirements.txt /code/requirements.txt
RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

# Копируем исходный код приложения
COPY ./app /code/app

# Копируем файл wait-for-it.sh
COPY ./wait-for-it.sh /code/wait-for-it.sh

# Делаем wait-for-it.sh исполнимым
RUN chmod +x /code/wait-for-it.sh

# Открываем порт 8000 для FastAPI
EXPOSE 8000

# Запуск FastAPI приложения через Uvicorn, после того как RabbitMQ будет доступен
CMD ["./wait-for-it.sh", "auth-rabbitmq:5672", "--", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]

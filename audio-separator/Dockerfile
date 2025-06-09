# Используем базовый образ с Python
FROM python:3.10-slim

RUN apt-get update && \
    apt-get install -y ffmpeg git && \
    pip install --upgrade pip

COPY app /app

WORKDIR /app

RUN pip install -r requirements.txt && \
    pip install demucs

CMD [ "python", "main.py" ]

FROM python:3.10-slim

RUN apt-get update && apt-get install -y \
    ffmpeg \
    libsndfile1 \
    && apt-get clean

WORKDIR /app

COPY /app/requirements.mlmodel.txt /app/

RUN pip install --no-cache-dir -r requirements.mlmodel.txt

COPY /app/requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

COPY /app/ /app/

CMD [ "python",  "server.py" ]
FROM python:3.10-slim

RUN apt-get update && \
    apt-get install -y ffmpeg
    
WORKDIR /app

COPY app/requirements.mlmodel.txt /app/

RUN pip install --no-cache-dir --default-timeout=1000 -r requirements.mlmodel.txt

COPY app/requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

COPY app/ /app/

RUN mkdir -p temp

CMD [ "python", "main.py" ]

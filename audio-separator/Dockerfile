FROM python:3.10-slim

RUN apt-get update && \
    apt-get install -y ffmpeg
    
WORKDIR /app

RUN pip install demucs

COPY app/requirements.txt /app/

RUN pip install -r requirements.txt

COPY app/ /app/

RUN mkdir -p temp

CMD [ "python", "main.py" ]

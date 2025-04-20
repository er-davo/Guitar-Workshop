FROM python:3.9

WORKDIR /audio-analyzer

COPY app /audio-analyzer

RUN apt-get update && apt-get install -y ffmpeg

RUN pip install grpcio grpcio-tools librosa supabase pydub

CMD [ "python", "main.py" ]
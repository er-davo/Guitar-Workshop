FROM python:3.9

WORKDIR /audio-analyzer

COPY app /audio-analyzer

RUN pip install grpcio grpcio-tools librosa supabase

CMD [ "python", "main.py" ]
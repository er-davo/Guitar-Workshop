FROM python:3.9

WORKDIR /audio-anlyzer

COPY app /audio-anlyzer

RUN pip install grpcio grpcio-tools

CMD [ "python", "main.py" ]
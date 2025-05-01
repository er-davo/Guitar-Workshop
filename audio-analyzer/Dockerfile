FROM python:3.9

WORKDIR /audio-analyzer

COPY app /audio-analyzer

RUN apt-get update && apt-get install -y ffmpeg

# RUN pip install -e .

RUN pip install --prefer-binary grpcio grpcio-tools librosa supabase resampy yt-dlp ffmpeg-python

CMD [ "python", "main.py" ]
import io
import librosa
from basic_pitch.inference import predict, Model
from basic_pitch import ICASSP_2022_MODEL_PATH


import tempfile
import soundfile as sf

class Analyzer:
    def __init__(self, sample_rate: int = 22050):
        self.sample_rate = sample_rate
        self.model = Model(ICASSP_2022_MODEL_PATH)

    def analyze(self, audio_bytes: bytes):
        y, sr = librosa.load(io.BytesIO(audio_bytes), sr=self.sample_rate, mono=True)

        with tempfile.NamedTemporaryFile(suffix=".wav", delete=True) as temp_wav:
            sf.write(temp_wav.name, y, sr)
            _, _, note_events = predict(temp_wav.name, self.model)

        return note_events



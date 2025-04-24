from files_access import download_file, delete_file, m4a_to_wav
from log import logger
from scipy.signal import butter, lfilter

import grpc
import librosa
import numpy as np

import audio_pb2_grpc
import audio_pb2


def butter_bandpass(lowcut, highcut, fs, order=5):
    nyq = 0.5 * fs
    low = lowcut / nyq
    high = highcut / nyq
    b, a = butter(order, [low, high], btype='band')
    return b, a

def bandpass_filter(data, lowcut, highcut, fs, order=5):
    b, a = butter_bandpass(lowcut, highcut, fs, order=order)
    y = lfilter(b, a, data)
    return y

class AudioAnalyzerServicer(audio_pb2_grpc.AudioAnalyzerServicer):
    def ProcessAudio(self, request: audio_pb2.AudioRequest, context):
        fmin = 80           # Игнорировать частоты ниже 80 Гц (басс)
        fmax = 1200         # Игнорировать частоты выше 1.2 кГц (вокал)
        n_fft = 4096        # Точность для низких нот (размер FFT окна)
        hop_length = 512    # Шаг для точности

        try:
            # Загрузка с подавлением низких частот (бас, барабаны)
            local_path = download_file(request.audio_path)
            logger.info(f"file {request.audio_path} downloaded")

            if request.audio_path[-4:] == ".m4a":
                request.audio_path = request.audio_path[:-3] + "wav"
                local_path = m4a_to_wav(request.audio_path)
                logger.info("converted from .m4a to wav")

            y, sr = librosa.load(local_path, sr=44100)
            y = librosa.util.normalize(y)
            if len(y) < hop_length:
                raise ValueError("Аудио слишком короткое для анализа")
            logger.info("loaded to librosa")

            delete_file(request.audio_path, del_supabase=False) # На время тестов
            
            # Уменьшение баса и барабанов
            y_filtered = bandpass_filter(y, fmin, fmax, sr)
            # y_filtered = librosa.effects.preemphasis(y=y, coef=0.97)
            y_filtered = librosa.effects.harmonic(y_filtered, margin=8)

            y_filtered = librosa.effects.trim(y_filtered, top_db=20)[0]  # Обрезка тишины
            # y_filtered = librosa.decompose.nn_filter(y_filtered)     # Подавление шумов

            
            f0, voiced_flag, _ = librosa.pyin(
                y_filtered,
                fmin=fmin,
                fmax=fmax,
                sr=sr,
                frame_length=2048,
                hop_length=hop_length,
                fill_na=np.nan
            )
            
            chroma = librosa.feature.chroma_cqt(
                y=y_filtered,
                sr=sr,
                hop_length=hop_length,
                bins_per_octave=36, # точность для гитары
                threshold=0.2
            )

            times = librosa.times_like(
                f0,
                sr=sr,
                hop_length=hop_length,
                n_fft=n_fft
            )

            results = []
            for i in range(len(f0)):
                event = audio_pb2.AudioEvent()
                time = times[i]
                event.time = float(time)

                if np.isnan(f0[i]):
                    continue
            
                solo_note_hz = f0[i]
                note = librosa.hz_to_note(solo_note_hz)

                chroma_frame = chroma[:, i]
                for j in range(12):
                    if chroma_frame[j] > 0.3:
                        chroma_note = librosa.midi_to_note(j + 36)[:-1]
                        event.chroma_notes.append(chroma_note)
                event.pitch = float(solo_note_hz)
                event.main_note = note[:-1]
                event.octave = int(note[-1])
            
                results.append(event)


            return audio_pb2.AudioResponse(
                note_features=results
            )
        
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error: {str(e)}")
            return audio_pb2.AudioResponse()

# на потом
# y_filtered = librosa.decompose.nn_filter(y_filtered)     # Подавление шумов

# Подавление негитарных частот
# y_guitar = librosa.effects.harmonic(y_filtered, margin=8)

# Разделение гармоник (гитара/вокал) и перкуссии (барабаны)
# y_harm, y_perc = librosa.decompose.hpss(y_filtered, margin=3.0) <- error
# if len(y_harm) < hop_length:
#     raise ValueError("Audio is too short for analysis")
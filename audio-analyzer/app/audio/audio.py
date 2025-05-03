from log import logger
from scipy.signal import butter, lfilter

import audio_pb2, audio_pb2_grpc
import storage
import youtube
import grpc
import librosa
import numpy as np
import os

test_files = [
    "nothing-else-matters.wav",
    "chords.wav",
]

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
        fmin = 70           # Минимальная частота анализа (для стандартного строя и Drop D)
        fmax = 1300         # Максимальная частота (увеличено для вокала)
        n_fft = 4096        # Размер FFT окна
        hop_length = 512    # Шаг анализа
        min_duration = 0.5  # Минимальная длительность аудио

        try:
            if request.type == audio_pb2.FILE:
                local_path = storage.download_file(request.audio_path)
                logger.info(f"File {request.audio_path} downloaded")
            elif request.type == audio_pb2.YOUTUBE:
                local_path = youtube.download_audio(request.audio_path)
                logger.info(f"File {local_path} downloaded")

            # Загрузка аудио с прогрессивной декодировкой
            y, sr = librosa.load(local_path, sr=44100, res_type='kaiser_best')
            y = librosa.util.normalize(y)

            if request.type == audio_pb2.FILE:
                storage.delete_file(request.audio_path, del_supabase=False if request.audio_path in test_files else True)
            elif request.type == audio_pb2.YOUTUBE:
                os.remove(local_path)
            
            # Проверка длительности аудио
            duration = len(y) / sr
            if duration < min_duration:
                raise ValueError(f"Audio too short ({duration:.2f}s < {min_duration}s)")

            if len(y) < hop_length:
                raise ValueError("Audio buffer too small for analysis")

            logger.info(f"Loaded audio: {duration:.2f} seconds, {sr} Hz")

            # Применение фильтров
            y_filtered = bandpass_filter(y, fmin, fmax, sr)
            # y_filtered = librosa.effects.harmonic(y_filtered, margin=8)
            y_filtered = librosa.effects.trim(y_filtered, top_db=30)[0]

            # HPSS разделение
            D = librosa.stft(y_filtered, n_fft=n_fft)
            H, P = librosa.decompose.hpss(D, margin=8.0)
            y_harm = librosa.istft(H)

            y_harm = librosa.effects.harmonic(
                y_filtered, 
                margin=6,
                kernel_size=9
            )

            # Pitch detection
            f0, voiced_flag, _ = librosa.pyin(
                y_harm,
                fmin=fmin,
                fmax=fmax,
                sr=sr,
                frame_length=n_fft,
                hop_length=hop_length,
                fill_na=0.0,  # Замена NaN на 0
                center=False
            )

            onset_frames = librosa.onset.onset_detect(
                y=y_harm,
                sr=sr,
                hop_length=hop_length,
                units='time'
            )

            # Хромаграмма
            chroma = librosa.feature.chroma_cqt(
                y=y_harm,
                sr=sr,
                hop_length=hop_length,
                bins_per_octave=36,
                threshold=0.2
            )

            # Временные метки
            times = librosa.times_like(f0, sr=sr, hop_length=hop_length)

            # Обработка результатов
            results = []
            for i in range(len(f0)):
                if f0[i] <= 0:  # Пропуск невалидных значений
                    continue

                event = audio_pb2.AudioEvent()
                event.time = float(times[i])
                event.pitch = float(f0[i])
                
                # Определение ноты
                note = librosa.hz_to_note(f0[i])
                event.main_note = note[:-1].replace("♯", "#")
                # event.main_note = event.main_note.replace("♯", "#")
                event.octave = int(note[-1])
                
                # Хроматические ноты
                chroma_frame = chroma[:, i]
                for j in range(12):
                    if chroma_frame[j] > 0.3:
                        midi_number = 36 + j + (event.octave * 12)
                        chroma_note = librosa.midi_to_note(midi_number)[:-1].replace("♯", "#")
                        event.chroma_notes.append(chroma_note)
                
                results.append(event)

            logger.info(f"Processed {len(results)} audio events")
            return audio_pb2.AudioResponse(note_features=results)

        except Exception as e:
            logger.error(f"Error processing {request.audio_path}: {str(e)}", exc_info=True)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Processing failed: {str(e)}")
            return audio_pb2.AudioResponse()
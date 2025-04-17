from files_access import download_file, delete_file

import grpc
import librosa
import numpy as np

import audio_pb2_grpc
import audio_pb2


class AudioAnalyzerServicer(audio_pb2_grpc.AudioAnalyzerServicer):
    def ProcessAudio(self, request : audio_pb2.AudioRequest, context):
        fmin = 80           # Игнорировать частоты ниже 80 Гц (басс)
        fmax = 1000         # Игнорировать частоты выше 1 кГц (вокал)
        n_fft = 4096        # Точность для низких нот (размер FFT окна)
        hop_length = 512    # Шаг для точности

        try:
            # Загрузка с подавлением низких частот (бас, барабаны)
            local_path = download_file(request.audio_path)
            y, sr = librosa.load(local_path, sr=44100)
            delete_file(request.audio_path)

            # Уменьшение баса и барабанов
            y_filtered = librosa.effects.preemphasis(y=y, coef=0.97)

            # Разделение гармоник (гитара/вокал) и перкуссии (барабаны)
            y_harm, y_perc = librosa.decompose.hpss(y_filtered, margin=3.0)

            f0, _, _ = librosa.pyin(
                y_harm,
                fmin=fmin,
                fmax=fmax,
                sr=sr,
                hop_length=hop_length
            )

            valid_mask = ~np.isnan(f0)
            pitches = f0[valid_mask].tolist()
            
            times = (np.arange(len(f0)) * hop_length / sr)[valid_mask].tolist()
            
            chroma = librosa.feature.chroma_cqt(
                y=y_harm,
                sr=sr,
                hop_length=hop_length
            )
            chroma_bytes = chroma.astype(np.float64).tobytes()

            tempo = librosa.beat.beat_track(y=y_harm, sr=sr)

            return audio_pb2.AudioResponse(
                pitches=pitches,
                times=times,
                chromagram=chroma_bytes,
                tempo=tempo,
                sr=sr,
                hop_length=hop_length
            )
        
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error: {str(e)}")
            return audio_pb2.AudioResponse()
from files_access import download_file, delete_file

import librosa
import numpy

import audio_pb2_grpc
import audio_pb2


class AudioAnalyzerServicer(audio_pb2_grpc.AudioAnalyzerServicer):
    def ProcessAudio(self, request : audio_pb2.AudioRequest, context):
        fmin = 80           # Игнорировать частоты ниже 80 Гц
        fmax = 1000         # Игнорировать частоты выше 1 кГц
        n_fft = 4096        # Точность для низких нот
        hop_length = 512    # Шаг для точности

        local_path = download_file(request.audio_path)

        y, sr = librosa.load(local_path, sr=None)

        delete_file(request.audio_path)

        pitches, magnitudes = librosa.piptrack(
            y=y,
            sr=sr,
            fmin=fmin,
            fmax=fmax,
            n_fft=n_fft,
            hop_length=hop_length
        )

        notes = []

        for t in range(pitches.shape[1]):
            frame_pitches = pitches[:, t]
            frame_magnitudes = magnitudes[:, t]

            if len(frame_magnitudes) > 0:
                max_idx = numpy.argmax(frame_magnitudes)
                freq = frame_pitches[max_idx]

                if freq > 0:
                    note = librosa.hz_to_note(freq)
                    notes.append(note)
        
        # for test
        print(notes)

        return audio_pb2.AudioResponse(notes=notes)
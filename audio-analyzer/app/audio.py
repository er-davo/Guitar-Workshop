import audio_pb2_grpc
import audio_pb2

class AudioAnalyzerServicer(audio_pb2_grpc.AudioAnalyzerServicer):
    def ProcessAudio(self, request, context):
        # TODO with librosa
        notes = ["C4", "G4", "A3"]
        return audio_pb2.AudioResponse(notes=notes)
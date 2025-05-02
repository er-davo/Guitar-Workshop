from concurrent import futures

import audio
import audio_pb2, audio_pb2_grpc
import grpc
import os


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    audio_pb2_grpc.add_AudioAnalyzerServicer_to_server(
        audio.AudioAnalyzerServicer(), server
    )
    
    server.add_insecure_port(f"0.0.0.0:{os.getenv('PORT')}")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
from concurrent import futures

import audio
import grpc


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    audio.add_AudioAnalyzerServicer_to_server(
        audio.AudioAnalyzerServicer(), server
    )
    
    server.add_insecure_port("0.0.0.0:50051")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()
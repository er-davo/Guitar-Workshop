from concurrent import futures
from service.service import AudioSeparatorService
import grpc
import separator_pb2_grpc

def run_server(port: str):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    separator_pb2_grpc.add_AudioSeparatorServicer_to_server(AudioSeparatorService("temp"), server)

    server.add_insecure_port(f"[::]:{port}")
    server.start()
    print(f"gRPC server is running on port {port}...")
    server.wait_for_termination()

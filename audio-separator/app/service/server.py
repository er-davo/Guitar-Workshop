from concurrent import futures
from service.service import AudioSeparatorService
import grpc
import separator_pb2_grpc

MAX_MESSAGE_LENGTH = 100 * 1024 * 1024  # 100 MB

def run_server(port: str):
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=4),
            options=[
            ('grpc.max_send_message_length', MAX_MESSAGE_LENGTH),
            ('grpc.max_receive_message_length', MAX_MESSAGE_LENGTH),
        ]
    )
    separator_pb2_grpc.add_AudioSeparatorServicer_to_server(AudioSeparatorService("temp"), server)

    server.add_insecure_port(f"[::]:{port}")
    server.start()
    print(f"gRPC server is running on port {port}...")
    server.wait_for_termination()

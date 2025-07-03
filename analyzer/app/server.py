import logging
import os
import grpc
from concurrent import futures
from analyzer_model import Analyzer
import note_analyzer_pb2
import note_analyzer_pb2_grpc

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s"
)

class NoteAnalyzerService(note_analyzer_pb2_grpc.NoteAnalyzerServicer):
    def __init__(self):
        self.analyzer = Analyzer()
        logging.info("NoteAnalyzerService initialized")

    def Analyze(self, request: note_analyzer_pb2.OAFRequest, context):
        file_name = request.audio_data.file_name
        logging.info(f"Received analysis request for file: {file_name}")

        try:
            notes = self.analyzer.analyze_bytes(
                audio_bytes=request.audio_data.audio_bytes,
                file_name=file_name
            )
            logging.info(f"Analysis completed successfully - {len(notes)} notes found")

            response = note_analyzer_pb2.OAFResponse()
            for note in notes:
                response.notes.add(
                    start_seconds=note.start_time,
                    midi_pitch=note.midi_pitch,
                    velocity=note.velocity,
                    duration_seconds=note.duration
                )
            return response

        except Exception as e:
            logging.exception(f"Error while analyzing file {file_name}: {e}")
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return note_analyzer_pb2.OAFResponse()

def serve():
    port = os.getenv("PORT")

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    note_analyzer_pb2_grpc.add_NoteAnalyzerServicer_to_server(NoteAnalyzerService(), server)
    server.add_insecure_port(f'[::]:{port}')

    logging.info(f"NoteAnalyzer gRPC server is running on port {port}")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()

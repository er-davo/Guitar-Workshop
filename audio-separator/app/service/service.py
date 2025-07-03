import grpc
import logging
from separator.core import Separator
import separator_pb2, separator_pb2_grpc

logger = logging.getLogger(__name__)

class AudioSeparatorService(separator_pb2_grpc.AudioSeparatorServicer):
    def __init__(self, path_to_temp_dir: str):
        super().__init__()
        self.temp_dir = path_to_temp_dir
        self.separator = Separator(self.temp_dir)
    
    def SeparateAudio(self, request: separator_pb2.SeparateRequest, context):
        file_name = request.audio_data.file_name
        logger.info(f"Received separation request for file: {file_name}")

        audio_bytes = request.audio_data.audio_bytes

        try:
            stems = self.separator.separate_audio_bytes(file_name, audio_bytes, cleanup=True)
        except Exception as e:
            logger.error(f"Separation failed for file {file_name}: {e}", exc_info=True)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return separator_pb2.SeparateResponse()
        
        response = separator_pb2.SeparateResponse()

        for stem_name, (fname, audio_bytes) in stems.items():
            logger.debug(f"Attaching stem: {stem_name} ({fname}, {len(audio_bytes)} bytes)")
            separated_audio_data = separator_pb2.AudioFileData(
                file_name=fname,
                audio_bytes=audio_bytes
            )
            response.stems[stem_name].CopyFrom(separated_audio_data)

        logger.info(f"Successfully separated {len(stems)} stems for file: {file_name}")
        return response

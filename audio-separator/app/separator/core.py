import logging
from typing import Dict, Tuple
import subprocess
import os
import shutil

logger = logging.getLogger(__name__)

class Separator:
    def __init__(self, path_to_temp_dir: str):
        self.temp_dir_path = path_to_temp_dir
        os.makedirs(self.temp_dir_path, exist_ok=True)
        logger.debug(f"Initialized Separator with temp dir: {self.temp_dir_path}")

    def separate_audio_bytes(
        self, 
        file_name: str, 
        audio_bytes: bytes, 
        cleanup: bool = True
    ) -> Dict[str, Tuple[str, bytes]]:
        logger.debug(f"Received audio size: {len(audio_bytes)} bytes")
        
        safe_file_name = file_name.replace(" ", "_")
        temp_file_path = os.path.join(self.temp_dir_path, safe_file_name)

        logger.debug(f"Saving audio to temp file: {temp_file_path}")
        with open(temp_file_path, "wb") as f:
            f.write(audio_bytes)

        output_dir = os.path.join(self.temp_dir_path, f"separated_{safe_file_name}")
        logger.debug(f"Output directory set to: {output_dir}")

        # try:
        #     result = subprocess.run(
        #         ["demucs", temp_file_path, "--out", output_dir, "--device", "cpu"],
        #         check=True,
        #         capture_output=True,
        #         text=True
        #     )
        #     logger.debug(f"Demucs stdout:\n{result.stdout}")
        #     if result.stderr:
        #         logger.warning(f"Demucs stderr:\n{result.stderr}")
        # except subprocess.CalledProcessError as e:
        #     logger.error(f"Demucs failed:\n{e.stderr}")
        #     raise RuntimeError("Demucs separation failed") from e

        result = subprocess.run(
            ["demucs", temp_file_path, "--out", output_dir],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        if result.returncode != 0:
            logger.error("Demucs failed with exit code %s", result.returncode)
            logger.error("STDOUT:\n%s", result.stdout)
            logger.error("STDERR:\n%s", result.stderr)  # Вот тут будет настоящая ошибка!
            raise RuntimeError("Demucs separation failed")

        subdirs = [
            d for d in os.listdir(output_dir)
            if os.path.isdir(os.path.join(output_dir, d))
        ]
        if not subdirs:
            raise RuntimeError(f"No model directory found in {output_dir}")
        model_name = subdirs[0]

        track_name = os.path.splitext(safe_file_name)[0]
        result_dir = os.path.join(output_dir, model_name, track_name)
        logger.debug(f"Looking for results in: {result_dir}")

        if not os.path.isdir(result_dir):
            raise RuntimeError(f"Result directory not found: {result_dir}")

        stems: Dict[str, Tuple[str, bytes]] = {}

        for fname in os.listdir(result_dir):
            stem_name = os.path.splitext(fname)[0]
            file_path = os.path.join(result_dir, fname)

            with open(file_path, "rb") as f:
                stems[stem_name] = (fname, f.read())

            logger.debug(f"Loaded stem: {stem_name} from {file_path}")

        if cleanup:
            logger.debug("Cleaning up temporary files...")
            try:
                os.remove(temp_file_path)
                shutil.rmtree(output_dir, ignore_errors=True)
            except Exception as e:
                logger.warning(f"Cleanup failed: {e}")

        logger.info(f"Separation complete. Extracted stems: {list(stems.keys())}")
        return stems

from typing import Dict, Tuple
import subprocess
import os
import shutil

class Separator:
    def __init__(self, path_to_temp_dir: str):
        self.temp_dir_path = path_to_temp_dir

    def separate_audio_bytes(
        self, 
        file_name: str, 
        audio_bytes: bytes, 
        cleanup: bool = True
    ) -> Dict[str, Tuple[str, bytes]]:

        temp_file_path = os.path.join(self.temp_dir_path, file_name)
        with open(temp_file_path, "wb") as f:
            f.write(audio_bytes)

        output_dir = os.path.join(self.temp_dir_path, f"separated_{file_name}")
        subprocess.run(
            [
                "demucs",
                temp_file_path,
                "--out", output_dir,
            ],
            check=True
        )

        model_name = os.listdir(output_dir)[0]
        track_name = os.path.splitext(file_name)[0]
        result_dir = os.path.join(output_dir, model_name, track_name)

        if not os.path.isdir(result_dir):
            raise RuntimeError(f"Result directory not found: {result_dir}")

        stems: Dict[str, Tuple[str, bytes]] = {}

        for fname in os.listdir(result_dir):
            stem_name = os.path.splitext(fname)[0]
            file_path = os.path.join(result_dir, fname)

            with open(file_path, "rb") as f:
                stems[stem_name] = (fname, f.read())

        if cleanup:
            os.remove(temp_file_path)
            shutil.rmtree(output_dir, ignore_errors=True)

        return stems

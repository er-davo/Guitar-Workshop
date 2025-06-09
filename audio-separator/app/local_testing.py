import sys
from pathlib import Path

from separator.core import Separator

def main():
    if len(sys.argv) != 2:
        print("Usage: python main.py path/to/audio.wav")
        return

    input_path = Path(sys.argv[1])
    if not input_path.exists():
        print(f"File {input_path} not found.")
        return

    with open(input_path, "rb") as f:
        audio_bytes = f.read()

    separator = Separator("temp")
    stems = separator.separate_audio_bytes(input_path.name, audio_bytes, cleanup=False)

    output_dir = Path("output_stems")
    output_dir.mkdir(exist_ok=True)

    for stem_name, (fname, data) in stems.items():
        out_path = output_dir / fname
        with open(out_path, "wb") as f:
            f.write(data)
        print(f"Saved: {out_path}")

if __name__ == "__main__":
    main()

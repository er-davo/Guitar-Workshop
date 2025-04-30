from supabase import create_client, Client
import os

SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("ACCESS_KEY")

LOCAL_STORAGE_PATH = "temp/"

supabase : Client = create_client(SUPABASE_URL, SUPABASE_KEY)

def upload_file(file, bucket_name: str = "audio-bucket") -> None:
    supabase.storage.from_(bucket_name).upload(file)

def download_file(file_path: str, bucket_name: str = "audio-bucket") -> str:
    res = supabase.storage.from_(bucket_name).download(file_path)
    with open(LOCAL_STORAGE_PATH + file_path, "wb") as f:
        f.write(res)
    print(f"File {file_path} downloaded")
    return "temp/" + file_path

def delete_file(file_path: str, bucket_name: str = "audio-bucket", del_supabase: bool = True ) -> None:
    try:
        if del_supabase:
            supabase.storage.from_(bucket_name).remove([file_path])
        os.remove(LOCAL_STORAGE_PATH + file_path)
    except Exception as e:
        print(f"Error deleting file {file_path}: {str(e)}")
        raise

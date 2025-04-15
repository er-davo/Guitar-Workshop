from supabase import create_client, Client
import os

SUPABASE_URL = os.getenv("SUPABASE_URL")
SUPABASE_KEY = os.getenv("ACCESS_KEY")

supabase : Client = create_client(SUPABASE_URL, SUPABASE_KEY)

def download_file(file_path : str, bucket_name : str = "audio-bucket") -> str:
    res = supabase.storage.from_(bucket_name).download(file_path)
    with open("temp/" + file_path, "wb") as f:
        f.write(res)
    print(f"File {file_path} downloaded")
    return "temp/" + file_path

def delete_file(file_path : str, bucket_name : str = "audio-bucket"):
    try:
        res = supabase.storage.from_(bucket_name).remove([file_path])
        os.remove("temp/" + file_path)
    except Exception as e:
        print(f"Error deleting file {file_path}: {str(e)}")
        raise

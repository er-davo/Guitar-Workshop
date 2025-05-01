import yt_dlp
import os

def download_audio(url, output_format="wav") -> str:
    """Downloads audio file from YouTube, return path to file"""
    ydl_opts = {
        'format': 'bestaudio/best',
        'outtmpl': 'temp/%(title)s.%(ext)s',  # Папка для сохранения
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'wav',
            'preferredquality': '192',        # Битрейт (для MP3)
        }],
        'quiet': True,                        # Убрать лишние логи
        'noplaylist': True,  # Skip playlists
        'extract_flat': False,
    }

    try:
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info = ydl.extract_info(url, download=True)
            filename = ydl.prepare_filename(info).replace(".webm", f".{output_format}").replace(".m4a", f".{output_format}")
            return filename  # Путь к скачанному файлу
    except Exception as e:
        return None
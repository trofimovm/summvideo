import os
import math
from fastapi import FastAPI, UploadFile, File, Form, HTTPException, Request
from fastapi.responses import JSONResponse, HTMLResponse, FileResponse, RedirectResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from moviepy import VideoFileClip
from pydub import AudioSegment
import tempfile
from openai import OpenAI
from datetime import datetime

app = FastAPI()

# Подключаем папку для статических файлов (например, CSS или изображения)
app.mount("/static", StaticFiles(directory="static"), name="static")

# Настраиваем шаблоны для рендеринга HTML страниц
templates = Jinja2Templates(directory="templates")

# Константы
PART_SIZE = 524288  # 512 KB
MAX_RETRIES = 3
MAX_FILE_SIZE = 500 * 1024 * 1024  # 500 MB

LOG_DIR = "/var/log/summvideo"
LOG_FILE = os.path.join(LOG_DIR, "log.txt")

def ensure_log_directory_exists():
    """Создаем директорию для логов, если её нет."""
    if not os.path.exists(LOG_DIR):
        os.makedirs(LOG_DIR)

def write_log(video_filename, prompt, transcription, summary):
    """Запись данных в лог."""
    ensure_log_directory_exists()

    with open(LOG_FILE, "a") as log_file:
        log_file.write(f"Дата и время: {datetime.now()}\n")
        log_file.write(f"Название видео: {video_filename}\n")
        log_file.write(f"Промт: {prompt}\n")
        log_file.write(f"Транскрибация: {transcription}\n")
        log_file.write(f"Саммари: {summary}\n")
        log_file.write("-" * 50 + "\n")

def extract_and_convert_audio(video_file, audio_file):
    """Функция для извлечения и конвертации аудио."""
    video = VideoFileClip(video_file)
    video.audio.write_audiofile(audio_file)
    
    audio = AudioSegment.from_wav(audio_file)
    mp3_file = audio_file.replace(".wav", ".mp3")
    audio.export(mp3_file, format="mp3", bitrate="32k")
    return mp3_file

def transcribe_audio(api_key, audio_file):
    """Функция для транскрибации аудио."""
    client = OpenAI(api_key=api_key)

    with open(audio_file, "rb") as file:
        transcript = client.audio.transcriptions.create(
            model="whisper-1",
            file=file,
            language="ru"
        )
    
    return transcript

def summarize_meeting_with_custom_prompt(api_key, transcript, prompt):
    """Использование пользовательского промта для генерации саммари."""
    client = OpenAI(api_key=api_key)

    response = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {"role": "system", "content": prompt},
            {"role": "user", "content": transcript}
        ],
        temperature=0,
        stream=True
    )

    collected_messages = []
    for chunk in response:
        chunk_message = chunk.choices[0].delta.content
        collected_messages.append(chunk_message)

    full_reply_content = ''.join([m for m in collected_messages if m is not None])
    return full_reply_content

@app.get("/", response_class=HTMLResponse)
async def read_root(request: Request):
    """Отображение HTML страницы приложения Vue."""
    # Запасной вариант - старый шаблон, если Vue.js сборка отсутствует
    return templates.TemplateResponse("upload_video.html", {"request": request})

@app.get("/index.html", response_class=HTMLResponse)
async def read_vue_index(request: Request):
    """Перенаправление на главную страницу."""
    return RedirectResponse(url="/")

@app.post("/upload_video/")
async def upload_video(file: UploadFile = File(...), prompt: str = Form(...)):
    """Обработчик для загрузки видео и выполнения обработки."""
    # Проверка размера файла
    if file.size > MAX_FILE_SIZE:
        raise HTTPException(status_code=400, detail="Размер файла превышает 500 МБ.")
    
    # Сохранение видеофайла во временное хранилище
    video_file = tempfile.NamedTemporaryFile(delete=False, suffix='.mp4')
    try:
        # Чтение и сохранение файла
        with open(video_file.name, "wb") as f:
            content = await file.read()  # Загружаем содержимое файла
            f.write(content)
        
        # Пути для временных аудиофайлов
        audio_file = tempfile.NamedTemporaryFile(delete=False, suffix='.wav').name
        
        # Извлечение аудио и его конвертация
        mp3_file = extract_and_convert_audio(video_file.name, audio_file)
        
        # Транскрибация аудио
        api_key = os.getenv("OPENAI_API_KEY")
        if not api_key:
            raise HTTPException(status_code=500, detail="API ключ OpenAI не найден.")
        
        transcription = transcribe_audio(api_key, mp3_file)
        transcription_text = transcription.text
        
        # Генерация саммари на основе переданного промта
        summary = summarize_meeting_with_custom_prompt(api_key, transcription_text, prompt)

        # Логирование
        write_log(file.filename, prompt, transcription_text, summary)

        # Включаем транскрипцию в ответ
        return JSONResponse(content={"summary": summary, "transcription": transcription_text})

    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

    finally:
        # Удаление временных файлов
        os.remove(video_file.name)
        os.remove(audio_file)
        if os.path.exists(mp3_file):
            os.remove(mp3_file)
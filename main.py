import os
import math
from fastapi import FastAPI, UploadFile, File, Form, HTTPException, Request
from fastapi.responses import JSONResponse, HTMLResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from moviepy.editor import VideoFileClip
from pydub import AudioSegment
import tempfile
from openai import OpenAI
import logging
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

# Настраиваем директорию для логов
LOG_DIR = "/tmp/summvideo"
LOG_FILE = os.path.join(LOG_DIR, "transcriptions_summary.log")

# Создаем директорию для логов, если она не существует
if not os.path.exists(LOG_DIR):
    try:
        os.makedirs(LOG_DIR)
        print(f"Директория {LOG_DIR} создана.")
    except Exception as e:
        print(f"Ошибка при создании директории {LOG_DIR}: {str(e)}")

# Настраиваем логирование
logging.basicConfig(
    filename=LOG_FILE,
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s",
)

def log_transcription(video_filename, prompt, transcription_text, summary):
    """Записываем информацию о транскрипции и саммари в лог-файл."""
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    log_message = (
        f"Дата: {timestamp}\n"
        f"Название видео: {video_filename}\n"
        f"Промт: {prompt}\n"
        f"Текст транскрипции: {transcription_text}\n"
        f"Саммари: {summary}\n"
        "------------------------\n"
    )
    try:
        logging.info(log_message)
        print("Информация записана в лог.")
    except Exception as e:
        print(f"Ошибка при записи в лог-файл: {str(e)}")

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
        model="gpt-4o",
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
    """Отображение HTML страницы для загрузки видео."""
    return templates.TemplateResponse("upload_video.html", {"request": request})

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

        # Логируем транскрипцию и саммари
        log_transcription(file.filename, prompt, transcription_text, summary)

        return JSONResponse(content={"summary": summary})

    except Exception as e:
        return JSONResponse(status_code=500, content={"error": str(e)})

    finally:
        # Удаление временных файлов
        os.remove(video_file.name)
        os.remove(audio_file)
        if os.path.exists(mp3_file):
            os.remove(mp3_file)
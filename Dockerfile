# Указываем базовый образ с Python
FROM python:3.10-slim

# Устанавливаем зависимости системы (включая ffmpeg для работы с видео и аудио)
RUN apt-get update && apt-get install -y ffmpeg

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /usr/src/app

# Копируем файл зависимостей и устанавливаем их
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Копируем все файлы проекта в контейнер
COPY . .

# Указываем переменные окружения
ENV OPENAI_API_KEY=ваш_ключ_openai

# Команда для запуска FastAPI приложения через Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
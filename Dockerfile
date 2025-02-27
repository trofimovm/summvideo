# Build backend directly without Vue.js for now
FROM python:3.10-slim

# Install system dependencies (including ffmpeg for video and audio processing)
RUN apt-get update && apt-get install -y ffmpeg

# Set working directory inside the container
WORKDIR /usr/src/app

# Copy dependency file and install dependencies
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy all project files to the container
COPY . .

# Environment variables
ENV OPENAI_API_KEY=ваш_ключ_openai

# Command to run the FastAPI application through Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
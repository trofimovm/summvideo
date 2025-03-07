#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}SummVideo - Запуск проекта${NC}"

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo "Docker не установлен. Пожалуйста, установите Docker и Docker Compose."
    exit 1
fi

# Проверка наличия .env файла
if [ ! -f .env ]; then
    echo -e "${YELLOW}Файл .env не найден. Создаю базовый .env файл...${NC}"
    echo "OPENAI_API_KEY=your_openai_api_key" > .env
    echo -e "${YELLOW}Пожалуйста, отредактируйте файл .env и добавьте ваш OPENAI_API_KEY${NC}"
fi

# Проверка наличия директории static/vue
if [ ! -d "static/vue" ]; then
    echo -e "${YELLOW}Директория static/vue не найдена. Необходимо собрать фронтенд...${NC}"
    
    # Проверка наличия Node.js и npm
    if command -v npm &> /dev/null; then
        echo -e "${GREEN}Сборка фронтенда...${NC}"
        cd frontend
        npm install
        npm run build
        cd ..
        echo -e "${GREEN}Фронтенд успешно собран!${NC}"
    else
        echo "Node.js и npm не установлены, но требуются для сборки фронтенда."
        echo "Установите Node.js и npm, затем запустите скрипт снова."
        echo "Или выполните следующие команды вручную:"
        echo "cd frontend && npm install && npm run build"
        exit 1
    fi
fi

# Запуск Docker Compose
echo -e "${GREEN}Запуск контейнеров Docker...${NC}"
docker compose up -d

# Проверка статуса контейнеров
echo -e "${GREEN}Проверка статуса контейнеров...${NC}"
docker compose ps

echo -e "${GREEN}SummVideo запущен и доступен по адресу: http://localhost:8000${NC}"
echo -e "${YELLOW}Для просмотра логов используйте: docker compose logs -f${NC}"
echo -e "${YELLOW}Для остановки используйте: docker compose down${NC}" 
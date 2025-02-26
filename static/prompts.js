// Глобальная переменная для хранения шаблонов
let promptTemplates = [];

// Загрузка шаблонов промптов из JSON файла
async function loadPromptTemplates() {
    try {
        const response = await fetch('/static/prompts.json');
        if (!response.ok) {
            throw new Error('Не удалось загрузить шаблоны промптов');
        }
        const data = await response.json();
        promptTemplates = data.templates;
        
        // Заполняем селект шаблонами
        populateTemplateSelect();
        
        // Применяем первый шаблон по умолчанию
        applyTemplate();
    } catch (error) {
        console.error('Ошибка при загрузке шаблонов:', error);
    }
}

// Заполнение селекта шаблонами из загруженного JSON
function populateTemplateSelect() {
    const templateSelect = document.getElementById('template-select');
    
    // Очищаем текущие опции
    templateSelect.innerHTML = '';
    
    // Добавляем опции из шаблонов
    promptTemplates.forEach(template => {
        const option = document.createElement('option');
        option.value = template.id;
        option.textContent = template.name;
        templateSelect.appendChild(option);
    });
    
    // Добавляем опцию "Свой промт"
    const customOption = document.createElement('option');
    customOption.value = 'custom_prompt';
    customOption.textContent = 'Свой промт';
    templateSelect.appendChild(customOption);
}

// Применение выбранного шаблона
function applyTemplate() {
    const promptTextarea = document.getElementById('prompt');
    const templateSelect = document.getElementById('template-select');
    const selectedTemplateId = templateSelect.value;
    
    if (selectedTemplateId === 'custom_prompt') {
        // Для "своего промта" просто очищаем поле
        promptTextarea.value = '';
    } else {
        // Ищем выбранный шаблон в массиве
        const selectedTemplate = promptTemplates.find(template => template.id === selectedTemplateId);
        if (selectedTemplate) {
            promptTextarea.value = selectedTemplate.text;
        }
    }
}

// Загрузка видео
async function uploadVideo(event) {
    event.preventDefault();

    const form = document.querySelector('form');
    const formData = new FormData(form);
    const submitButton = document.querySelector('button[type="submit"]');

    // Отключаем кнопку отправки
    submitButton.disabled = true;
    submitButton.textContent = 'Обработка...';

    // Показываем прелоудер и сообщение
    document.getElementById('loader').style.display = 'block';
    document.getElementById('processing-message').style.display = 'block';
    document.getElementById('result').innerHTML = '';
    document.getElementById('transcription-container').style.display = 'none';

    try {
        const response = await fetch('/upload_video/', {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('Ошибка загрузки видео');
        }

        const result = await response.json();
        
        // Преобразуем Markdown в HTML
        const md = window.markdownit({
            html: true, // разрешаем HTML внутри Markdown
            linkify: true, // автоформатируем ссылки
            typographer: true // используем типограф
        });
        const htmlContent = md.render(result.summary);

        document.getElementById('result').innerHTML = htmlContent;
        
        // Отображаем транскрипцию
        if (result.transcription) {
            document.getElementById('transcription-text').innerText = result.transcription;
            document.getElementById('transcription-container').style.display = 'block';
        }
    } catch (error) {
        document.getElementById('result').innerText = 'Ошибка: ' + error.message;
    } finally {
        // Скрываем прелоудер и сообщение
        document.getElementById('loader').style.display = 'none';
        document.getElementById('processing-message').style.display = 'none';

        // Включаем кнопку отправки обратно
        submitButton.disabled = false;
        submitButton.textContent = 'Отправить';
    }
}

// Функция для переключения видимости транскрипции
function toggleTranscription() {
    const content = document.getElementById('transcription-content');
    const button = document.getElementById('transcription-toggle');
    
    if (content.style.display === 'none' || content.style.display === '') {
        content.style.display = 'block';
        button.textContent = 'Скрыть транскрипцию';
    } else {
        content.style.display = 'none';
        button.textContent = 'Показать транскрипцию';
    }
}

// Загружаем шаблоны при загрузке страницы
window.onload = function() {
    loadPromptTemplates();
};
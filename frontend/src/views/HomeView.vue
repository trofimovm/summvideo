<template>
  <div class="container">
    <header class="header">
      <div class="logo">
        <span class="logo-emoji">🎙️</span>
        <h1>Умные резюме видео</h1>
      </div>
      <p class="description">
        SummVideo превращает любое видео в текст и создаёт информативное резюме за считанные минуты.
        Экономьте до 90% времени на обработке видеоматериалов с помощью передовых AI-технологий.
      </p>
    </header>
    
    <main class="main">
      <section class="content-section">
        <TelegramAuth v-if="!isAuthenticated" />
        
        <article class="benefits-section" v-if="!isAuthenticated">
          <h2>Преимущества SummVideo</h2>
          <div class="benefits-grid">
            <div class="benefit-card">
              <span class="benefit-icon">⏱️</span>
              <h3>Экономия времени</h3>
              <p>Получайте результат за минуты вместо часов ручной работы</p>
            </div>
            <div class="benefit-card">
              <span class="benefit-icon">🎯</span>
              <h3>Высокая точность</h3>
              <p>Транскрипция с точностью до 98% при хорошем качестве аудио</p>
            </div>
            <div class="benefit-card">
              <span class="benefit-icon">📊</span>
              <h3>Умные резюме</h3>
              <p>Структурированные конспекты, выделяющие главное из контента</p>
            </div>
            <div class="benefit-card">
              <span class="benefit-icon">🧩</span>
              <h3>Гибкая настройка</h3>
              <p>Выбор шаблонов резюме под конкретные типы контента</p>
            </div>
            <div class="benefit-card">
              <span class="benefit-icon">🔄</span>
              <h3>Автоматизация</h3>
              <p>Полностью автоматический процесс от загрузки видео до получения готового резюме</p>
            </div>
            <div class="benefit-card">
              <span class="benefit-icon">🔍</span>
              <h3>Полная транскрипция</h3>
              <p>Получайте не только резюме, но и полный текст видео для дальнейшего использования</p>
            </div>
          </div>
        </article>
      
        <template v-else>
          <UploadForm v-if="!isProcessing && !summary" />
          
          <ProcessingIndicator v-if="isProcessing" />
          
          <div class="results" v-if="summary">
            <h2>Результат анализа видео</h2>
            <ResultDisplay :summary="summary" />
            
            <div class="actions">
              <button class="btn btn-primary" @click="resetForm">
                <span class="icon">🔄</span> Обработать еще видео
              </button>
            </div>
            
            <TranscriptionViewer v-if="transcription" :transcription="transcription" />
          </div>
        </template>
      </section>
      
      <section class="faq-section" v-if="!isAuthenticated">
        <h2>Часто задаваемые вопросы</h2>
        <div class="faq-grid">
          <div class="faq-item">
            <h3>Какие форматы видео поддерживаются?</h3>
            <p>SummVideo поддерживает все популярные форматы видео (MP4, AVI, MOV, WMV, MKV). Максимальный размер файла — 500 МБ.</p>
          </div>
          <div class="faq-item">
            <h3>Как настроить формат резюме?</h3>
            <p>Выберите готовый шаблон (деловые встречи, образовательные видео) или создайте свой запрос с указанием структуры.</p>
          </div>
          <div class="faq-item">
            <h3>Как долго обрабатывается видео?</h3>
            <p>В среднем 3-5 минут для видео продолжительностью 1 час — значительно быстрее ручной обработки.</p>
          </div>
          <div class="faq-item">
            <h3>Сколько стоит использование сервиса?</h3>
            <p>Первые транскрипции и резюме — бесплатно! Зарегистрируйтесь через Telegram для доступа ко всем функциям.</p>
          </div>
        </div>
      </section>
    </main>
    
    <footer class="footer">
      <p>© {{ new Date().getFullYear() }} SummVideo — умные резюме видео</p>
    </footer>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import UploadForm from '@/components/UploadForm.vue';
import ProcessingIndicator from '@/components/ProcessingIndicator.vue';
import ResultDisplay from '@/components/ResultDisplay.vue';
import TranscriptionViewer from '@/components/TranscriptionViewer.vue';
import TelegramAuth from '@/components/TelegramAuth.vue';

export default {
  name: 'HomeView',
  components: {
    UploadForm,
    ProcessingIndicator,
    ResultDisplay,
    TranscriptionViewer,
    TelegramAuth
  },
  setup() {
    const store = useStore();

    const resetForm = () => {
      store.commit('CLEAR_RESULTS');
    };

    return {
      isProcessing: computed(() => store.getters.isProcessing),
      summary: computed(() => store.getters.getSummary),
      transcription: computed(() => store.getters.getTranscription),
      error: computed(() => store.getters.getError),
      isAuthenticated: computed(() => store.getters['auth/isAuthenticated']),
      resetForm
    };
  }
}
</script>

<style scoped>
.container {
  width: 100%;
  max-width: 1000px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  text-align: center;
  margin-bottom: 2.5rem;
  padding-bottom: 1.8rem;
  border-bottom: 1px solid var(--gray-color);
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.logo-emoji {
  font-size: 2.5rem;
  margin-right: 1rem;
}

.logo h1 {
  font-size: 2.2rem;
  margin: 0;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.3;
  max-width: 800px;
}

.description {
  max-width: 800px;
  margin: 0 auto;
  color: #4b5563;
  font-size: 1.15rem;
  line-height: 1.6;
}

.main {
  flex: 1;
}

.content-section {
  background-color: white;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 2.5rem;
  margin-bottom: 2.5rem;
}

.benefits-section {
  margin-bottom: 3rem;
}

.benefits-section h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--dark-color);
  font-size: 1.8rem;
}

.benefits-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
  gap: 1.5rem;
}

.benefit-card {
  padding: 1.5rem;
  border-radius: var(--border-radius);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  text-align: center;
}

.benefit-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.benefit-icon {
  font-size: 2.5rem;
  display: block;
  margin-bottom: 1rem;
}

.benefit-card h3 {
  margin: 0 0 0.8rem;
  font-size: 1.3rem;
  color: var(--primary-color);
}

.benefit-card p {
  color: #4b5563;
  margin: 0;
  line-height: 1.5;
}

.faq-section {
  background-color: white;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 2.5rem;
  margin-bottom: 2.5rem;
}

.faq-section h2 {
  text-align: center;
  margin-bottom: 2rem;
  color: var(--dark-color);
  font-size: 1.8rem;
}

.faq-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.faq-item {
  padding: 1.5rem;
  border-radius: var(--border-radius);
  background-color: #f9fafb;
  border-left: 4px solid var(--primary-color);
}

.faq-item h3 {
  margin: 0 0 0.8rem;
  font-size: 1.2rem;
  color: var(--dark-color);
}

.faq-item p {
  color: #4b5563;
  margin: 0;
  line-height: 1.5;
}

.results h2 {
  text-align: center;
  margin-bottom: 1.8rem;
  color: var(--dark-color);
  font-size: 1.8rem;
}

.actions {
  display: flex;
  justify-content: center;
  margin: 2.5rem 0;
}

.btn {
  display: inline-flex;
  align-items: center;
  padding: 0.8rem 1.8rem;
  border-radius: var(--border-radius);
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  border: none;
  font-size: 1.1rem;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: #2779bd;
  transform: translateY(-3px);
  box-shadow: 0 5px 15px rgba(0, 110, 220, 0.2);
}

.icon {
  margin-right: 0.5rem;
}

.footer {
  text-align: center;
  padding: 2rem 0;
  color: #4b5563;
  font-size: 0.95rem;
  border-top: 1px solid var(--gray-color);
  margin-top: 1.5rem;
}

@media (max-width: 768px) {
  .header {
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
  }
  
  .content-section, .faq-section {
    padding: 1.8rem;
  }

  .logo-emoji {
    font-size: 2rem;
  }

  .logo h1 {
    font-size: 1.8rem;
  }

  .description {
    font-size: 1rem;
    padding: 0 1rem;
  }
  
  .benefits-grid, .faq-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .content-section, .faq-section {
    padding: 1.5rem;
  }
  
  .benefit-card, .faq-item {
    padding: 1.2rem;
  }
  
  .logo h1 {
    font-size: 1.6rem;
  }
}
</style>
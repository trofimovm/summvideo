<template>
  <div class="container">
    <header class="header">
      <div class="logo">
        <span class="logo-emoji">🎙️</span>
        <h1>SummVideo</h1>
      </div>
      <p class="description">
        Сервис для создания кратких саммари из видеоконтента.
        С помощью AI-технологий он выделяет главные моменты из любого видео,
        экономя ваше время и предоставляя только самую важную информацию.
      </p>
    </header>
    
    <main class="main">
      <section class="content-section">
        <UploadForm v-if="!isProcessing && !summary" />
        
        <ProcessingIndicator v-if="isProcessing" />
        
        <div class="results" v-if="summary">
          <h2>Результат анализа</h2>
          <ResultDisplay :summary="summary" />
          
          <div class="actions">
            <button class="btn btn-primary" @click="resetForm">
              <span class="icon">🔄</span> Обработать еще видео
            </button>
          </div>
          
          <TranscriptionViewer v-if="transcription" :transcription="transcription" />
        </div>
      </section>
    </main>
    
    <footer class="footer">
      <p>© {{ new Date().getFullYear() }} SummVideo - AI-помощник для анализа видео</p>
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

export default {
  name: 'HomeView',
  components: {
    UploadForm,
    ProcessingIndicator,
    ResultDisplay,
    TranscriptionViewer
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
      resetForm
    };
  }
}
</script>

<style scoped>
.container {
  width: 100%;
  max-width: 900px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  text-align: center;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--gray-color);
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.logo-emoji {
  font-size: 2.5rem;
  margin-right: 0.75rem;
}

.logo h1 {
  font-size: 2.5rem;
  margin: 0;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.description {
  max-width: 700px;
  margin: 0 auto;
  color: #64748b;
  font-size: 1.1rem;
}

.main {
  flex: 1;
}

.content-section {
  background-color: white;
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 2rem;
  margin-bottom: 2rem;
}

.results h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: var(--dark-color);
}

.actions {
  display: flex;
  justify-content: center;
  margin: 2rem 0;
}

.btn {
  display: inline-flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  border-radius: var(--border-radius);
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
  border: none;
  font-size: 1rem;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background-color: #2779bd;
  transform: translateY(-2px);
}

.icon {
  margin-right: 0.5rem;
}

.footer {
  text-align: center;
  padding: 1.5rem 0;
  color: #64748b;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .content-section {
    padding: 1.5rem;
  }

  .logo-emoji {
    font-size: 2rem;
  }

  .logo h1 {
    font-size: 2rem;
  }

  .description {
    font-size: 1rem;
    padding: 0 1rem;
  }
}

@media (max-width: 480px) {
  .content-section {
    padding: 1rem;
  }
}
</style>
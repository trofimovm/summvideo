<template>
  <div class="transcription-container">
    <button class="toggle-btn" @click="toggleTranscription">
      <span class="icon">{{ isVisible ? '📜' : '📄' }}</span>
      {{ isVisible ? 'Скрыть транскрипцию' : 'Показать полную транскрипцию' }}
    </button>
    
    <transition name="slide-fade">
      <div class="transcription-panel" v-if="isVisible">
        <div class="transcription-header">
          <h3>Полная транскрипция</h3>
          <button class="copy-btn" @click="copyToClipboard" :class="{ copied: isCopied }">
            <span class="copy-icon">{{ isCopied ? '✓' : '📋' }}</span> 
            {{ isCopied ? 'Скопировано!' : 'Копировать' }}
          </button>
        </div>
        <div class="transcription-content">
          <div class="transcription-text" ref="transcriptionEl">{{ transcription }}</div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'TranscriptionViewer',
  props: {
    transcription: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const isVisible = ref(false);
    const isCopied = ref(false);
    const transcriptionEl = ref(null);

    const toggleTranscription = () => {
      isVisible.value = !isVisible.value;
    };
    
    const copyToClipboard = () => {
      if (!transcriptionEl.value) return;
      
      navigator.clipboard.writeText(props.transcription)
        .then(() => {
          isCopied.value = true;
          setTimeout(() => {
            isCopied.value = false;
          }, 2000);
        })
        .catch(err => {
          console.error('Не удалось скопировать текст: ', err);
        });
    };

    return {
      isVisible,
      toggleTranscription,
      transcriptionEl,
      copyToClipboard,
      isCopied
    };
  }
}
</script>

<style scoped>
.transcription-container {
  margin-top: 2rem;
  width: 100%;
}

.toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 0.75rem;
  background-color: var(--light-color);
  color: var(--dark-color);
  border: 1px solid var(--gray-color);
  border-radius: var(--border-radius);
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
}

.toggle-btn:hover {
  background-color: #e2e8f0;
}

.icon {
  margin-right: 0.5rem;
  font-size: 1.25rem;
}

.transcription-panel {
  margin-top: 1rem;
  background-color: white;
  border-radius: var(--border-radius);
  border: 1px solid var(--gray-color);
  overflow: hidden;
}

.transcription-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: #f8fafc;
  border-bottom: 1px solid var(--gray-color);
}

.transcription-header h3 {
  margin: 0;
  font-size: 1.1rem;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  background-color: var(--light-color);
  color: var(--dark-color);
  border: 1px solid var(--gray-color);
  border-radius: var(--border-radius);
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: var(--transition);
}

.copy-btn:hover {
  background-color: #e2e8f0;
}

.copy-btn.copied {
  background-color: #d1fae5;
  color: var(--success-color);
  border-color: var(--success-color);
}

.copy-icon {
  margin-right: 0.5rem;
}

.transcription-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 1.5rem;
}

.transcription-text {
  font-size: 0.95rem;
  line-height: 1.6;
  white-space: pre-wrap;
  color: #4a5568;
}

/* Transitions */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

@media (max-width: 768px) {
  .transcription-content {
    padding: 1rem;
  }
  
  .transcription-text {
    font-size: 0.9rem;
  }
}
</style>
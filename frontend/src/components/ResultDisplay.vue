<template>
  <div class="result-container">
    <div class="result-header">
      <div class="divider">
        <span class="divider-icon">✨</span>
      </div>
    </div>
    <div class="result-content-wrapper">
      <div class="result-actions">
        <button class="copy-btn" @click="copyToClipboard" :class="{ copied: isCopied }">
          <span class="copy-icon">{{ isCopied ? '✓' : '📋' }}</span> 
          {{ isCopied ? 'Скопировано!' : 'Копировать' }}
        </button>
      </div>
      <div class="result-content" v-html="formattedSummary" ref="contentEl"></div>
    </div>
  </div>
</template>

<script>
import { computed, ref } from 'vue';
import MarkdownIt from 'markdown-it';

export default {
  name: 'ResultDisplay',
  props: {
    summary: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const md = new MarkdownIt({
      html: true,
      linkify: true,
      typographer: true
    });

    const contentEl = ref(null);
    const isCopied = ref(false);

    const formattedSummary = computed(() => {
      return props.summary ? md.render(props.summary) : '';
    });

    const copyToClipboard = () => {
      if (!contentEl.value) return;
      
      // Получаем только текст из HTML-содержимого
      const tempEl = document.createElement('div');
      tempEl.innerHTML = contentEl.value.innerHTML;
      const textToCopy = tempEl.textContent || tempEl.innerText || '';
      
      navigator.clipboard.writeText(textToCopy)
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
      formattedSummary,
      contentEl,
      copyToClipboard,
      isCopied
    };
  }
}
</script>

<style scoped>
.result-container {
  width: 100%;
  margin-top: 1rem;
}

.result-header {
  margin-bottom: 2rem;
  text-align: center;
  position: relative;
}

.divider {
  display: flex;
  align-items: center;
  justify-content: center;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--gray-color);
}

.divider::before {
  margin-right: 1rem;
}

.divider::after {
  margin-left: 1rem;
}

.divider-icon {
  font-size: 1.5rem;
}

.result-content-wrapper {
  position: relative;
}

.result-actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 0.75rem;
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

.result-content {
  background-color: white;
  border-radius: var(--border-radius);
  padding: 2rem;
  box-shadow: var(--box-shadow);
  overflow-wrap: break-word;
}

.result-content :deep(h1), 
.result-content :deep(h2), 
.result-content :deep(h3) {
  color: var(--dark-color);
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-weight: 600;
}

.result-content :deep(h1) {
  font-size: 1.75rem;
}

.result-content :deep(h2) {
  font-size: 1.5rem;
}

.result-content :deep(h3) {
  font-size: 1.25rem;
}

.result-content :deep(p) {
  margin-bottom: 1rem;
  line-height: 1.6;
}

.result-content :deep(ul), 
.result-content :deep(ol) {
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

.result-content :deep(li) {
  margin-bottom: 0.5rem;
}

.result-content :deep(blockquote) {
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
  margin-left: 0;
  margin-right: 0;
  color: #4a5568;
}

.result-content :deep(code) {
  background-color: #f1f5f9;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
}

.result-content :deep(pre) {
  background-color: #1e293b;
  color: #f8fafc;
  padding: 1rem;
  border-radius: var(--border-radius);
  overflow-x: auto;
  margin-bottom: 1rem;
}

.result-content :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  color: inherit;
}

.result-content :deep(a) {
  color: var(--primary-color);
  text-decoration: none;
}

.result-content :deep(a:hover) {
  text-decoration: underline;
}

.result-content :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 0 auto;
  border-radius: var(--border-radius);
}

@media (max-width: 768px) {
  .result-content {
    padding: 1.5rem;
  }
}

@media (max-width: 480px) {
  .result-content {
    padding: 1rem;
  }
  
  .divider::before,
  .divider::after {
    max-width: 60px;
  }
}
</style>
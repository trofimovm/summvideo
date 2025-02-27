<template>
  <div class="upload-form-container">
    <h2>Загрузите видео для анализа</h2>
    
    <form @submit.prevent="submitForm" class="upload-form">
      <div class="form-group">
        <label for="file" class="file-label">
          <div class="file-upload-area" :class="{ 'has-file': selectedFileName }">
            <div class="file-icon">
              <span v-if="!selectedFileName">📤</span>
              <span v-else>🎥</span>
            </div>
            <div class="file-info">
              <span v-if="!selectedFileName">Выберите или перетащите видеофайл</span>
              <span v-else>{{ selectedFileName }}</span>
              <small>Максимальный размер файла: 500 МБ</small>
            </div>
          </div>
          <input 
            type="file" 
            id="file" 
            ref="fileInput" 
            accept="video/*" 
            required 
            @change="handleFileChange"
            class="file-input"
          />
        </label>
      </div>

      <div class="form-group">
        <label for="template-select">Шаблон промта:</label>
        <select 
          id="template-select" 
          v-model="selectedTemplateId" 
          @change="handleTemplateChange"
          class="form-control"
        >
          <option 
            v-for="template in promptTemplates" 
            :key="template.id" 
            :value="template.id"
          >
            {{ template.name }}
          </option>
          <option value="custom_prompt">Свой промт</option>
        </select>
      </div>

      <div class="form-group">
        <label for="prompt">Промт для обработки:</label>
        <textarea 
          id="prompt" 
          v-model="promptText" 
          required
          class="form-control"
          rows="6"
        ></textarea>
      </div>

      <div class="form-actions">
        <button type="submit" class="btn btn-primary" :disabled="isProcessing || !selectedFile">
          <span class="icon">🚀</span> {{ isProcessing ? 'Обработка...' : 'Отправить' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useStore } from 'vuex';

export default {
  name: 'UploadForm',
  setup() {
    const store = useStore();
    const fileInput = ref(null);
    const selectedFile = ref(null);
    const selectedFileName = ref('');
    const selectedTemplateId = ref('');
    const promptText = ref('');
    
    const promptTemplates = computed(() => store.getters.getPromptTemplates);
    const isProcessing = computed(() => store.getters.isProcessing);

    onMounted(() => {
      store.dispatch('loadPromptTemplates').then(() => {
        // Select first template by default if available
        if (promptTemplates.value.length > 0) {
          selectedTemplateId.value = promptTemplates.value[0].id;
          applyTemplate();
        }
      });
    });
    
    // Watch for template selection changes
    watch(selectedTemplateId, () => {
      applyTemplate();
    });

    const handleFileChange = (event) => {
      const file = event.target.files[0];
      if (file) {
        selectedFile.value = file;
        selectedFileName.value = file.name;
      } else {
        selectedFile.value = null;
        selectedFileName.value = '';
      }
    };

    const applyTemplate = () => {
      if (selectedTemplateId.value === 'custom_prompt') {
        promptText.value = '';
        return;
      }
      
      const template = promptTemplates.value.find(
        template => template.id === selectedTemplateId.value
      );
      
      if (template) {
        promptText.value = template.text;
      }
    };

    const submitForm = () => {
      if (!selectedFile.value || !promptText.value) {
        return;
      }
      
      store.dispatch('processVideo', {
        file: selectedFile.value,
        prompt: promptText.value
      });
    };

    return {
      fileInput,
      selectedFile,
      selectedFileName,
      selectedTemplateId,
      promptText,
      promptTemplates,
      isProcessing,
      handleFileChange,
      handleTemplateChange: applyTemplate,
      submitForm
    };
  }
}
</script>

<style scoped>
.upload-form-container {
  width: 100%;
}

h2 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: var(--dark-color);
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: var(--dark-color);
}

.file-label {
  cursor: pointer;
  display: block;
}

.file-upload-area {
  border: 2px dashed #cbd5e0;
  border-radius: var(--border-radius);
  padding: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition);
  background-color: #f8fafc;
}

.file-upload-area:hover {
  border-color: var(--primary-color);
  background-color: #ebf4ff;
}

.file-upload-area.has-file {
  border-color: var(--success-color);
  background-color: #ebfff4;
}

.file-icon {
  font-size: 2.5rem;
  margin-right: 1rem;
}

.file-info {
  display: flex;
  flex-direction: column;
}

.file-info small {
  color: #718096;
  margin-top: 0.25rem;
}

.file-input {
  display: none;
}

.form-control {
  padding: 0.75rem;
  border: 1px solid #cbd5e0;
  border-radius: var(--border-radius);
  font-size: 1rem;
  transition: var(--transition);
}

.form-control:focus {
  border-color: var(--primary-color);
  outline: none;
  box-shadow: 0 0 0 3px rgba(52, 144, 220, 0.2);
}

textarea.form-control {
  resize: vertical;
  min-height: 120px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
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

.btn-primary:hover:not(:disabled) {
  background-color: #2779bd;
  transform: translateY(-2px);
}

.btn-primary:disabled {
  background-color: #a0aec0;
  cursor: not-allowed;
}

.icon {
  margin-right: 0.5rem;
}

@media (max-width: 768px) {
  .file-upload-area {
    padding: 1.5rem;
    flex-direction: column;
    text-align: center;
  }
  
  .file-icon {
    margin-right: 0;
    margin-bottom: 0.5rem;
  }
  
  .form-actions {
    justify-content: center;
  }
}
</style>
<template>
  <div class="processing-container">
    <div class="processing-content">
      <div class="spinner-container">
        <div class="spinner"></div>
      </div>
      <h3>Обработка видео</h3>
      <p class="processing-message">
        Ваше видео обрабатывается с помощью AI. Это может занять несколько минут в зависимости от размера файла.
      </p>
      <div class="processing-steps">
        <div class="step">
          <div class="step-icon">📤</div>
          <div class="step-desc">Извлечение аудио</div>
        </div>
        <div class="step">
          <div class="step-icon">🎙️</div>
          <div class="step-desc">Транскрибация текста</div>
        </div>
        <div class="step">
          <div class="step-icon">🧠</div>
          <div class="step-desc">Анализ содержания</div>
        </div>
        <div class="step">
          <div class="step-icon">📝</div>
          <div class="step-desc">Создание саммари</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProcessingIndicator'
}
</script>

<style scoped>
.processing-container {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 2rem 0;
}

.processing-content {
  text-align: center;
  max-width: 600px;
}

h3 {
  margin: 1rem 0;
  color: var(--dark-color);
}

.spinner-container {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
}

.spinner {
  width: 60px;
  height: 60px;
  border: 5px solid rgba(52, 144, 220, 0.2);
  border-radius: 50%;
  border-top-color: var(--primary-color);
  animation: spin 1.5s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.processing-message {
  color: #64748b;
  margin-bottom: 2rem;
}

.processing-steps {
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
  flex-wrap: wrap;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 25%;
  padding: 0 0.5rem;
  transition: var(--transition);
}

.step-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.step-desc {
  font-size: 0.875rem;
  color: #4a5568;
  text-align: center;
}

@media (max-width: 640px) {
  .processing-steps {
    flex-wrap: wrap;
  }
  
  .step {
    width: 50%;
    margin-bottom: 1.5rem;
  }
}
</style>
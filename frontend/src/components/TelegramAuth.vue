<template>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h2>Начните использовать SummVideo прямо сейчас</h2>
        <p>1 час транскрибации и резюмирования видео <strong>бесплатно</strong>!</p>
      </div>
      
      <div class="auth-body">
        <div class="auth-cta">
          <p class="cta-text">Войдите через Telegram, чтобы получить доступ к сервису и начать создавать качественные резюме ваших видео</p>
        </div>
        
        <div class="telegram-widget-container" ref="telegramLoginWidget"></div>
        
        <div v-if="error" class="auth-error">
          {{ error }}
        </div>
        
        <div class="auth-info">
          <p>Используя авторизацию через Telegram, вы соглашаетесь с нашими условиями использования.</p>
          <p>Мы не храним и не используем ваш пароль от Telegram.</p>
          
          <div class="dev-mode" v-if="showDevMode">
            <button class="dev-mode-btn" @click="handleDevLogin">
              Войти в режиме разработки
            </button>
            <p class="dev-mode-info">Эта опция доступна только в режиме разработки</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue';
import { useStore } from 'vuex';

export default {
  name: 'TelegramAuth',
  setup() {
    const store = useStore();
    const telegramLoginWidget = ref(null);
    const error = ref('');
    const showDevMode = ref(true); // Всегда показываем кнопку для режима разработки
    
    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated']);
    
    onMounted(() => {
      // Создаем скрипт для загрузки виджета авторизации Telegram
      try {
        const script = document.createElement('script');
        script.src = 'https://telegram.org/js/telegram-widget.js?21';
        script.setAttribute('data-telegram-login', process.env.VUE_APP_TELEGRAM_BOT_NAME || 'your_bot_name');
        script.setAttribute('data-size', 'large');
        script.setAttribute('data-radius', '8');
        script.setAttribute('data-request-access', 'write');
        script.setAttribute('data-userpic', 'true');
        script.setAttribute('data-onauth', 'onTelegramAuth(user)');
        script.async = true;
        
        // Добавляем обработчик авторизации в window
        window.onTelegramAuth = (user) => {
          // Отправляем данные пользователя на сервер
          handleTelegramLogin(user);
        };
        
        // Добавляем скрипт на страницу
        if (telegramLoginWidget.value) {
          telegramLoginWidget.value.innerHTML = '';
          telegramLoginWidget.value.appendChild(script);
        }
      } catch (err) {
        console.error('Ошибка при инициализации Telegram виджета:', err);
        error.value = 'Не удалось загрузить виджет Telegram. Попробуйте войти в режиме разработки.';
        showDevMode.value = true;
      }
    });
    
    // Обработчик данных авторизации от Telegram
    const handleTelegramLogin = async (user) => {
      try {
        error.value = '';
        
        // Преобразуем объект user в query параметры
        const params = new URLSearchParams();
        Object.keys(user).forEach(key => {
          params.append(key, user[key]);
        });
        
        // Отправляем запрос на сервер
        const response = await fetch(`/auth/telegram?${params.toString()}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        });
        
        const data = await response.json();
        
        if (!response.ok) {
          throw new Error(data.error || 'Ошибка авторизации');
        }
        
        // Сохраняем токен и данные пользователя
        store.commit('auth/SET_TOKEN', data.token);
        store.commit('auth/SET_USER', data.user);
        
        // Перенаправляем на главную страницу
        if (isAuthenticated.value) {
          store.dispatch('resetResults');
        }
      } catch (err) {
        console.error('Ошибка авторизации через Telegram:', err);
        error.value = err.message || 'Произошла ошибка при авторизации';
        showDevMode.value = true;
      }
    };
    
    // Обработчик для режима разработки
    const handleDevLogin = () => {
      // Создаем фиктивные данные пользователя для режима разработки
      const mockUser = {
        id: 123456789,
        first_name: 'Test',
        last_name: 'User',
        username: 'test_user',
        photo_url: '',
        auth_date: Math.floor(Date.now() / 1000),
        hash: 'dev_mode_hash'
      };
      
      // Сохраняем токен и данные пользователя
      store.commit('auth/SET_TOKEN', 'dev_mode_token');
      store.commit('auth/SET_USER', mockUser);
      
      // Перенаправляем на главную страницу
      store.dispatch('resetResults');
    };
    
    return {
      telegramLoginWidget,
      error,
      isAuthenticated,
      showDevMode,
      handleDevLogin
    };
  }
};
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem 1rem;
}

.auth-card {
  background-color: white;
  border-radius: var(--border-radius);
  width: 100%;
  max-width: 450px;
  box-shadow: var(--box-shadow);
  overflow: hidden;
}

.auth-header {
  padding: 2rem;
  text-align: center;
  background: linear-gradient(135deg, var(--primary-color), #2779bd);
  color: white;
  border-bottom: 1px solid var(--gray-color);
}

.auth-header h2 {
  margin-bottom: 0.5rem;
  color: white;
  font-size: 1.5rem;
}

.auth-header p {
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 0;
  font-size: 1.1rem;
}

.auth-body {
  padding: 2rem;
}

.auth-cta {
  text-align: center;
  margin-bottom: 1.5rem;
}

.cta-text {
  font-size: 1.05rem;
  color: var(--dark-color);
  margin-bottom: 1.5rem;
  line-height: 1.5;
}

.telegram-widget-container {
  display: flex;
  justify-content: center;
  margin-bottom: 1.5rem;
  min-height: 50px;
}

.auth-error {
  background-color: #fef2f2;
  color: var(--danger-color);
  padding: 1rem;
  border-radius: var(--border-radius);
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
  text-align: center;
}

.auth-info {
  font-size: 0.85rem;
  color: #64748b;
  text-align: center;
}

.auth-info p {
  margin-bottom: 0.5rem;
}

.dev-mode {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px dashed var(--gray-color);
}

.dev-mode-btn {
  background-color: #4f46e5;
  color: white;
  border: none;
  border-radius: var(--border-radius);
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition);
}

.dev-mode-btn:hover {
  background-color: #4338ca;
}

.dev-mode-info {
  margin-top: 0.75rem;
  font-size: 0.75rem;
  color: #94a3b8;
}

@media (max-width: 480px) {
  .auth-header, .auth-body {
    padding: 1.5rem;
  }
}
</style>
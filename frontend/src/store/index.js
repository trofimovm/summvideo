import { createStore } from 'vuex'
import ApiService from '../services/ApiService'

export default createStore({
  state: {
    promptTemplates: [],
    isProcessing: false,
    summary: '',
    transcription: '',
    error: null
  },
  getters: {
    getPromptTemplates: state => state.promptTemplates,
    isProcessing: state => state.isProcessing,
    getSummary: state => state.summary,
    getTranscription: state => state.transcription,
    getError: state => state.error
  },
  mutations: {
    SET_PROMPT_TEMPLATES(state, templates) {
      state.promptTemplates = templates;
    },
    SET_PROCESSING(state, status) {
      state.isProcessing = status;
    },
    SET_SUMMARY(state, summary) {
      state.summary = summary;
    },
    SET_TRANSCRIPTION(state, transcription) {
      state.transcription = transcription;
    },
    SET_ERROR(state, error) {
      state.error = error;
    },
    CLEAR_RESULTS(state) {
      state.summary = '';
      state.transcription = '';
      state.error = null;
    }
  },
  actions: {
    async loadPromptTemplates({ commit }) {
      try {
        const templates = await ApiService.getPromptTemplates();
        commit('SET_PROMPT_TEMPLATES', templates);
      } catch (error) {
        commit('SET_ERROR', 'Не удалось загрузить шаблоны промптов');
        console.error('Error loading prompt templates:', error);
      }
    },
    async processVideo({ commit }, { file, prompt }) {
      commit('CLEAR_RESULTS');
      commit('SET_PROCESSING', true);
      
      try {
        const result = await ApiService.uploadVideo(file, prompt);
        commit('SET_SUMMARY', result.summary);
        commit('SET_TRANSCRIPTION', result.transcription);
      } catch (error) {
        commit('SET_ERROR', 'Ошибка обработки видео: ' + error.message);
        console.error('Error processing video:', error);
      } finally {
        commit('SET_PROCESSING', false);
      }
    }
  }
})
import axios from 'axios';

// Base API URL - adjust as needed for production/development
const API_URL = '';

const ApiService = {
  /**
   * Get prompt templates from the JSON file
   */
  async getPromptTemplates() {
    try {
      const response = await axios.get(`${API_URL}/static/prompts.json`);
      return response.data.templates;
    } catch (error) {
      console.error('Error fetching prompt templates:', error);
      throw error;
    }
  },

  /**
   * Upload video file and prompt for processing
   * @param {File} file - The video file to upload
   * @param {String} prompt - The prompt text for processing
   * @param {String} token - JWT auth token
   */
  async uploadVideo(file, prompt, token) {
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('prompt', prompt);

      const headers = {
        'Content-Type': 'multipart/form-data'
      };
      
      // Добавляем заголовок авторизации, если токен предоставлен
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await axios.post(`${API_URL}/upload_video/`, formData, { headers });

      return response.data;
    } catch (error) {
      console.error('Error uploading video:', error);
      throw error;
    }
  }
};

export default ApiService;
const state = {
  token: localStorage.getItem('token') || '',
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  status: ''
};

const getters = {
  isAuthenticated: state => !!state.token,
  getUser: state => state.user,
  getStatus: state => state.status,
  getToken: state => state.token
};

const mutations = {
  SET_TOKEN(state, token) {
    state.token = token;
    localStorage.setItem('token', token);
  },
  SET_USER(state, user) {
    state.user = user;
    localStorage.setItem('user', JSON.stringify(user));
  },
  SET_STATUS(state, status) {
    state.status = status;
  },
  LOGOUT(state) {
    state.token = '';
    state.user = null;
    state.status = '';
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }
};

const actions = {
  // Действие для авторизации через Telegram
  telegramAuth({ commit }, authData) {
    return new Promise((resolve, reject) => {
      commit('SET_STATUS', 'loading');
      const params = new URLSearchParams();
      Object.keys(authData).forEach(key => {
        params.append(key, authData[key]);
      });

      fetch(`/auth/telegram?${params.toString()}`)
        .then(response => response.json())
        .then(data => {
          if (data.error) {
            commit('SET_STATUS', 'error');
            reject(data.error);
          } else {
            commit('SET_TOKEN', data.token);
            commit('SET_USER', data.user);
            commit('SET_STATUS', 'success');
            resolve(data);
          }
        })
        .catch(err => {
          commit('SET_STATUS', 'error');
          reject(err);
        });
    });
  },

  // Действие для проверки текущего статуса авторизации
  checkAuth({ commit, state }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        commit('LOGOUT');
        reject(new Error('No auth token'));
        return;
      }

      commit('SET_STATUS', 'loading');
      fetch('/profile', {
        headers: {
          'Authorization': `Bearer ${state.token}`
        }
      })
        .then(response => {
          if (!response.ok) {
            throw new Error('Token invalid');
          }
          return response.json();
        })
        .then(data => {
          commit('SET_USER', data.user);
          commit('SET_STATUS', 'success');
          resolve(data);
        })
        .catch(err => {
          commit('LOGOUT');
          reject(err);
        });
    });
  },

  // Действие для выхода из системы
  logout({ commit }) {
    return new Promise(resolve => {
      commit('LOGOUT');
      resolve();
    });
  }
};

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
};
import Vue from 'vue'
import Vuex from 'vuex'
import Axios from 'axios'

Vue.prototype.$http = Axios
const token = localStorage.getItem('token')
if (token) {
  Vue.prototype.$http.defaults.headers.common['Authorization'] = token
}

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    status: '',
    token: localStorage.getItem('token') || '',
    user: {}
  },
  mutations: {
    auth_request (state) {
      state.status = 'loading'
    },
    auth_success (state, token, user) {
      state.status = 'success'
      state.token = token
      state.user = user
    },
    auth_error (state) {
      state.status = 'error'
    },
    whoami_request (state) {
      state.status = 'loading'
    },
    whoami_success (state, token, user) {
      state.status = 'success'
      state.user = user
    },
    whoami_error (state) {
      state.status = 'error'
    },
    logout (state) {
      state.status = ''
      state.token = ''
    }
  },
  getters: {
    isLoggedIn: state => !!state.token,
    authStatus: state => state.status
  },
  actions: {
    whoami ({commit}) {
      return new Promise((resolve, reject) => {
        commit('auth_request')
        this.$http({ url: 'http://localhost:8084/whoami', method: 'GET' })
          .then(resp => {
            const user = resp.data.user
            localStorage.setItem('user', user)
            commit('whoami_success', user)
            resolve(resp)
          })
          .catch(err => {
            commit('whoami_error')
            localStorage.removeItem('token')
            reject(err)
          })
      })
    },
    tokenFromCode ({ commit }, code) {
      return new Promise((resolve, reject) => {
        this.$http({ url: 'http://localhost:3000/login', data: { code: code }, method: 'POST' }).then(resp => {
          const token = resp.data.token
          localStorage.setItem('token', token)
          this.$http.defaults.headers.common['Authorization'] = token
          commit('auth_success', token)
          resolve(resp)
        }).catch(err => {
          commit('auth_error')
          localStorage.removeItem('token')
          reject(err)
        })
      })
    }
  }
})

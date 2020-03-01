import Vue from 'vue'
import Vuex from 'vuex'
import VueAxios from 'vue-axios'
import VueAuthenticate from 'vue-authenticate'
import axios from 'axios'

Vue.use(Vuex)
Vue.use(VueAxios, axios)

var vueAuth = VueAuthenticate.factory(Vue.prototype.$http, {
  baseUrl: 'http://localhost:8083'
})

export default new Vuex.Store({
  state: {
    isAuthenticated: false
  },
  getters: {
    isAuthenticated () {
      return vueAuth.isAuthenticated()
    }
  },
  mutations: {
    isAuthenticated (state, payload) {
      state.isAuthenticated = payload.isAuthenticated
    }
  },
  actions: {
    login (context, payload) {
      vueAuth.login(payload.user, payload.requestOptions).then((response) => {
        context.commit('isAuthenticated', {
          isAuthenticated: vueAuth.isAuthenticated()
        })
      })
    }
  }
})

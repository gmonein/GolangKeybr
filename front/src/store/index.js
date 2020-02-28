import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    authenticated: false
  },
  mutations: {
    authenticate (state) {
      console.log('tata')
      state.authenticated = true
    },
    logout (state) {
      console.log('toto')
      state.authenticated = false
    }
  }
})

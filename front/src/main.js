// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import VueAxios from 'vue-axios'
import axios from 'axios'
import VueCookies from 'vue-cookies'
import router from './router'
import store from './store/index'
import userService from '@/_services/user.service'

Vue.use(VueCookies)

Vue.config.productionTip = false

var authenticated = false
/* eslint-disable no-new */
var vm = new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>',
  name: 'App',
  data: {
    authenticated: authenticated
  }
})
vm.userService = userService

Vue.use(VueAxios, axios)

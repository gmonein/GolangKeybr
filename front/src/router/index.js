import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginComponent from '@/components/login'
import SecureComponent from '@/components/secure'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      redirect: {
        name: 'login'
      }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginComponent
    },
    {
      path: '/secure',
      name: 'secure',
      component: SecureComponent
    }
  ]
})

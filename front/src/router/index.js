import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginPage from '@/components/LoginPage'
import HomePage from '@/components/HomePage'
import userService from '@/_services/user.service'

Vue.use(Router)

export const router = new Router({
  mode: 'history',
  routes: [
    { path: '/', component: HomePage },
    { path: '/login', component: LoginPage },

    { path: '*', redirect: '/' }
  ]
})

router.beforeEach((to, from, next) => {
  userService.loggedIn().then(function (e) {
    if (e !== false) {
      userService.setLogin(e['login'])
      router.app.$data.authenticated = true
      if (to.path === '/login') {
        next('/')
      } else {
        next()
      }
    } else {
      router.app.$data.authenticated = false
      next('/login')
    }
  }).catch(e => {
    console.log(e, 'catcehd')
  })
  next()
})

export default router

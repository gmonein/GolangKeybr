import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import HomePage from '@/components/HomePage'
// import userService from '@/_services/user.service'

Vue.use(Router)

function oauthCallbackHandler (to, from, next) {
  let code = to.query.code
  axios({ url: 'http://localhost:3000/login', data: {code: code}, method: 'GET' })
}

export const router = new Router({
  mode: 'history',
  routes: [
    { path: '/', component: HomePage },
    { path: '/oauth/marvin', beforeEnter: oauthCallbackHandler },
    { path: '*', redirect: '/' }
  ]
})

export default router

import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginPage from '@/components/LoginPage'
import HomePage from '@/components/HomePage'
// import userService from '@/_services/user.service'
import Oauth from '@/_services/oauth'

Vue.use(Router)

export const router = new Router({
  mode: 'history',
  routes: [
    { path: '/', component: HomePage, meta: { requireAuth: true } },
    { path: '/login', component: LoginPage },
    { path: '/oauth/marvin',
      beforeEnter: Oauth.oauthCallbackHandler },
    { path: '*', redirect: '/' }
  ]
})

router.beforeEach((to, from, next) => {
  if (to.meta.requireAuth && !localStorage.getItem('user')) {
    next('/login')
    return
  }
  next()
})

export default router

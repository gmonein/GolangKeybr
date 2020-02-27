import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginComponent from '@/components/login'
import SecureComponent from '@/components/secure'

Vue.use(Router)

let router = new Router({
  mode: 'history',
  routes: [
    { path: '/', name: 'home', component: LoginComponent },
    { path: '/login', name: 'login', component: LoginComponent },
    { path: '/secure', name: 'secure', component: SecureComponent, meta: { requiresAuth: true } }]
})

function userIsAuth () {
  return true
}

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth) && !userIsAuth()) {
    next({ name: 'login' })
    return
  }
  if (to.name === 'login' && userIsAuth()) {
    next({ name: 'secure' })
    return
  }
  next()
})

export default router

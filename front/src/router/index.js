import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginComponent from '@/components/login'
import SecureComponent from '@/components/secure'

Vue.use(Router)

let router = new Router({
  routes: [
    { path: '/', redirect: { name: 'login' } },
    { path: '/login', name: 'login', component: LoginComponent, meta: { respondToToken: true } },
    { path: '/secure', name: 'secure', component: SecureComponent, meta: { requiresAuth: true } }
  ]
})

function userIsAuth () {
  return false
}

router.beforeEach((to, from, next) => {
  console.log(to)
  console.log(from)
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!userIsAuth()) {
      next({ name: 'login' })
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router

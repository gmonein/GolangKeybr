import Vue from 'vue'
import Router from 'vue-router'
// import HelloWorld from '@/components/HelloWorld'
import LoginPage from '@/components/LoginPage'
import HomePage from '@/components/HomePage'
// import userService from '@/_services/user.service'

Vue.use(Router)

export const router = new Router({
  mode: 'history',
  routes: [
    { path: '/', component: HomePage },
    { path: '/login', component: LoginPage },
    { path: '/oauth/marvin',
      beforeEnter: (to, from, next) => {
        let code = to.query.code
        if (code) {
          fetch(`http://localhost:8082/oauth?code=${code}`)
            .then(e => {
              if (e.status === 200) {
                return e.json()
              } else {
                return new Promise(() => { return false })
              }
            }).then(e => {
              console.log(e)
              if (e === false) {
                next('/login')
              } else {
                const parseJwt = (token) => {
                  try {
                    return JSON.parse(atob(token.split('.')[1]))
                  } catch (e) {
                    return null
                  }
                }
                console.log(parseJwt(e.token))
                next('/')
              }
            })
        } else {
          next('/login')
        }
      }},
    { path: '*', redirect: '/' }
  ]
})

router.beforeEach((to, from, next) => {
  // userService.loggedIn().then(function (e) {
  //   if (e !== false) {
  //     userService.setLogin(e['login'])
  //     router.app.$data.authenticated = true
  //     if (to.path === '/login') {
  //       next('/')
  //     } else {
  //       next()
  //     }
  //   } else {
  //     router.app.$data.authenticated = false
  //     next('/login')
  //   }
  // }).catch(e => {
  //   console.log(e, 'catcehd')
  // })
  next()
})

export default router

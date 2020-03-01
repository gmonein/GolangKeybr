export default {
  oauthCallbackHandler
}

function oauthCallbackHandler (to, from, next) {
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
          localStorage.setItem('user', parseJwt(e.token).login)
          localStorage.setItem('token', e.token)
          console.log(parseJwt(e.token))
          next('/')
        }
      })
  } else {
    next('/login')
  }
}

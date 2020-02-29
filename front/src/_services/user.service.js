import config from '@/../config/index'

const userService = {
  login,
  loggedIn,
  logout,
  setLogin
}

function setLogin (login) {
  localStorage.setItem('user', login)
}

function login (username, password) {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  }

  return fetch(`${config.apiUrl}/users/authenticate`, requestOptions)
    .then(handleResponse)
    .then(user => {
      if (user.token) {
        localStorage.setItem('user', JSON.stringify(user))
      }

      return user
    })
}

function loggedIn () {
  console.log('do fetch')
  return fetch('http://localhost:8082/whoami', { credentials: 'include' })
    .then(e => {
      console.log(e)
      return e
    })
    .then(e => {
      if (e.status === 200) {
        return e.json()
      } else {
        return new Promise(() => { return false })
      }
    })
}

function logout () {
  return fetch('http://localhost:8082/logout', { credentials: 'include' }).then(e => {
    if (e.status === 200) {
      return e.json()
    } else {
      return new Promise(() => { return false })
    }
  }).catch(e => {
    return new Promise(() => { return false })
  })
}

function handleResponse (response) {
  return response.text().then(text => {
    const data = text && JSON.parse(text)
    if (!response.ok) {
      if (response.status === 401) {
        location.reload(true)
      }

      const error = (data && data.message) || response.statusText
      return Promise.reject(error)
    }

    return data
  })
}

export default userService

<template>
  <div id="app">
    <router-link v-if="authenticated" to="/login" v-on:click.native="logout()" replace>Logout</router-link>
    <a href="/secure">Logout</a>
    <router-view @authenticated="setAuthenticated" />
  </div>
</template>

<script>
export default {
  name: 'App',
  data () {
    return {
      authenticated: false
    }
  },
  methods: {
    setAuthenticated (status) {
      this.$store.authenticated = status
    },
    logout () {
      this.$store.authenticated = this.$store.commit('authenticate')
    }
  },
  mounted () {
    fetch('http://localhost:8082/whoami', {
      credentials: 'include'
    }).then(e => {
      e.status === 200 ? this.$store.commit('authenticate') : this.$store.commit('logout')
    })
  },
  watch: {
    $route (to, from) {
      fetch('http://localhost:8082/whoami', {
        credentials: 'include'
      }).then(e => {
        e.status === 200 ? this.$store.commit('authenticate') : this.$store.commit('logout')
      })
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>

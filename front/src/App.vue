<template>
  <div id="app">
    <router-link v-if="!authenticated" to="#" v-on:click.native="logout()" replace>Logout</router-link>
    <router-view/>
  </div>
</template>

<script>
export default {
  data () {
    return {
      authenticated: false
    }
  },
  methods: {
    setAuthenticated (status) {
      this.$data.authenticated = status
      this.authenticated = status
    },
    logout () {
      this.$root.userService.logout()
      console.log('toto')
      this.setAuthenticated(false)
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

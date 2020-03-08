<template>
  <div id="login">
    <div v-if="!authenticated">
      <h1>Login</h1>
      <button @click="authenticate()">Sign with intranet</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Login',
  methods: {
    authenticate: function () {
      window.location.href = 'https://api.intra.42.fr/oauth/authorize?client_id=08fdc53b395d24655e44ba228b44e2e52b4ed63e683e610a1798bbfa56a39aa6&redirect_uri=http%3A%2F%2Flocalhost%3A8083%2Foauth%2Fmarvin&response_type=code'
    },
    logout: function () {
      this.$store.dispatch('logout')
    }
  },
  computed: {
    isLoggedIn: function () { return this.$store.getters.isLoggedIn }
  }
}

</script>

<style scoped>
#login {
  width: 500px;
  border: 1px solid #CCCCCC;
  background-color: #FFFFFF;
  margin: auto;
  margin-top: 200px;
  padding: 20px;
}
</style>

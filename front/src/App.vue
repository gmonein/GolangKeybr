<template>
  <div id="app">
    <router-link to="#" v-on:click.native="logout()" replace>Logout</router-link>
    <router-view/>
  </div>
</template>

<script>
export default {
  methods: {
    logout () {
      localStorage.removeItem('user', null)
      this.$router.push('/login')
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

  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

* {
  border-color: #657b83;
}
a {
  color: #268bd2;
}
a:visited {
  color: #6c71c4;
}
body {
  background-color: #002b36;
  color: #839496;
}
html {
  background-color: #002b36;
}
input,
textarea {
  background-color: #073642;
  color: #839496;
}
blockquote,
pre {
  background-color: #073642;
  color: #839496;
}
* {
  border-color: #657b83;
}
a {
  color: #268bd2;
}
a:visited {
  color: #6c71c4;
}
body {
  background-color: #1d2021;
  color: #839496;
}
html {
  background-color: #1d2021;
}
input,
textarea {
  background-color: #073642;
  color: #839496;
}
blockquote,
pre {
  background-color: #073642;
  color: #839496;
}
  .black {
    color:   '#282828'
  }
  .red {
    color:     '#cc241d'
  }
  .green {
    color:   '#98971a'
  }
  .yellow {
    color:  '#d79921'
  }
  .blue {
    color:    '#458588'
  }
  .magenta {
    color: '#b16286'
  }
  .cyan {
    color:    '#689d6a'
  }

  .white {
    color:   '#a89984'
  }

</style>

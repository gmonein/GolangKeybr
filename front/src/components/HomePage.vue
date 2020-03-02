<template>
  <div id='citation'>
    <p><span v-for="letter in citationH" :key="letter.index" :class="letter.class">{{ letter.letter }}</span></p>
    <a href="#" v-bind:class="showMore" @click="showMore"> Click Me </a>
  </div>
</template>

<script>
export default {
  name: 'citation',
  data: function () {
    return {
      citationH: [],
      citation: '',
      index: 0
    }
  },
  getters: {
    signalSocket: state => state.signalSocket,
    dataSocket: state => state.dataSocket
  },
  mutations: {
    signalSocket (state, v) {
      state.signalSocket = v
    },
    dataSocket (state, v) {
      state.dataSocket = v
    }
  },
  methods: {
    showMore: function () {
      this.index = this.index + 1
    },

    initCit: function () {
      this.citationH = []
      for (let i = 0; i < this.citation.length; i++) {
        this.citationH.push({index: i, letter: this.citation[i], class: ''})
      }
      console.log(this.citationH)
    },
    addSpan: function () {
      this.citationBeforeIndex = this.citation.slice(0, this.index)
      this.currentIndexLetter = this.citation[this.index]
      this.citationAfterIndex = this.citation.slice(this.index + 1)
      console.log(this.citationAfterIndex)
    },
    connectType: function () {
    }
  },
  mounted () {
    console.log(this.$root.Websocket('ws://localhost:8082/type_ws'))
    let socket = new this.$root.Websocket('ws://localhost:8082/type_ws')
    socket.onopen = function (e) {
      alert('[open] Connection established')
      alert('Sending to server')
      socket.send('My name is John')
    }

    socket.onmessage = function (event) {
      alert(`[message] Data received from server: ${event.data}`)
    }

    socket.onclose = function (event) {
      if (event.wasClean) {
        alert(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`)
      } else {
        // e.g. server process killed or network down
        // event.code is usually 1006 in this case
        alert('[close] Connection died')
      }
    }

    this.citation = 'Please wait'
    console.log(localStorage.getItem('token'))
    fetch('http://localhost:8082/citation', {
      headers: {
        Authorization: localStorage.getItem('token')
      }
    }).then(e => {
      console.log(e)
      if (e.status === 200) {
        e.text().then(t => { this.citation = t })
      }
    })
  },
  watch: {
    citation: function (val) {
      this.initCit()
      this.addSpan()
    },
    index: function (val) {
      this.citationH[this.index - 1].class = 'typed'
      this.citationH[this.index].class = 'index'
    }
  }
}
</script>

<style scoped>
.index {
  background-color: #2f3e43
}

.typed {
  opacity: 70%
}

.error {
  color: red
}
</style>

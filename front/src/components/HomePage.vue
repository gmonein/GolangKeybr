<template>
  <div id='citation'>
    <p><span v-for="letter in citationH" :key="letter.index" :class="letter.class">{{ letter.letter }}</span></p>
    <input v-on:input="handleInput" placeholder="modifiez-moi">
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
      index: 0,
      typeSocket: {
        socket: null,
        connected: false,
        error: ''
      }
    }
  },
  methods: {
    handleInput: function (event) {
      if (event.inputType !== 'insertText') { return }
      if (!this.typeSocket.connected) { return }
      this.typeSocket.socket.send('1' + event.data)
    },
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
    }
  },
  mounted () {
    this.typeSocket.socket = new WebSocket('ws://localhost:8082/type_ws')
    this.typeSocket.socket.onopen = e => {
      this.typeSocket.socket.send(localStorage.getItem('token'))
    }
    this.typeSocket.socket.onmessage = e => {
      switch (e.data[0]) {
        case '1':
          this.index = this.index + 1
          break
        case '2':
          console.log('auth err')
          break
        case '3':
          console.log('auth succ')
          this.typeSocket.connected = true
          break
        case '4':
          this.citation = e.data.slice(1)
          break
        default:
          console.log(e.data)
      }
    }
    this.typeSocket.socket.onclose = e => {
      console.log(e)
    }
    this.citation = 'Please wait'
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

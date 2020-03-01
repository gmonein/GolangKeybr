<template>
  <body>
    <div>
      <p><span v-for="letter in citation">{{ letter }}</span></p>
      <p>{{ citationBeforeIndex }}<span class="index">{{ currentIndexLetter }}</span>{{ citationAfterIndex }}</p>
      <a href="#" v-bind:class="showMore" @click="showMore"> Click Me </a>
    </div>
  </body>
</template>

<script>
export default {
  name: 'HelloWorld',
  data () {
    return {
      citation: '',
      citationBeforeIndex: '',
      currentIndexLetter: '',
      citationAfterIndex: '',
      index: 0
    }
  },
  methods: {
    showMore: function () {
      this.index = this.index + 1
    },
    addSpan: function () {
      this.citationBeforeIndex = this.citation.slice(0, this.index)
      this.currentIndexLetter = this.citation[this.index]
      this.citationAfterIndex = this.citation.slice(this.index + 1)
      console.log(this.citationAfterIndex)
    }
  },
  mounted () {
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
      this.addSpan()
    },
    index: function (val) {
      this.addSpan()
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.index {
  background-color: #6c71c4;
}
</style>

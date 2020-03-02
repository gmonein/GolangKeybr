<template>
  <p><span v-for="letter in citationH" :key="letter.index" :class="letter.class">{{ letter.letter }}</span></p>
</template>

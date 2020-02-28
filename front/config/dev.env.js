'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

console.log('tata')
module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  externals: {
      // global app config object
      config: JSON.stringify({
          apiUrl: 'http://localhost:8082'
      })
})

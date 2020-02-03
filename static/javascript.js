fetch('/citation').then( e => {
console.log(e) ; e.blob() }).then( e => {
  console.log(e)
  document.getElementById('citation').html = e
})

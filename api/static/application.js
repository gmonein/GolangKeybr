let userInput = document.getElementById('userInput')
let userLoginButton = document.getElementById('userLoginButton')
let letterInput = document.getElementById('letterInput')
let usersTable = document.getElementById('usersTable')
let windowLocation = window.location.host
let socketPath = (window.location.protocol == "https:" ? "wss:/" : "ws:/") + "/" + window.location.host

let getCitation = () => {
fetch("/citation")
.then(response => response.text())
.then((response) => {
  let div = document.getElementById('citation')
  let explode = [...response]
  let index = -1
  div.innerHTML = explode.map( letter => {
    if (letter == "\n") {
      index += 1
      return `<div index=${index} class='letter space'> </div><br>`
    } else if (letter == " ") {
      index += 1
      return `<div index=${index} class='letter space'> </div>`
    } else {
      index += 1
      return `<div index=${index} class='letter'>${letter}</div>`
    }
  }).join('')
  letterInput.index = 0
  letterInput.maxIndex = explode.length
  letterInput.value = ''
  document.querySelector(`div.letter[index='0']`).classList.add('index')
})
.catch(err => console.log(err))
}

let openGetSocket = function() {
window.getSocket = new WebSocket(`${socketPath}/data_ws`)
window.getSocket.onopen = () => { window.getSocket.send(JSON.stringify({ name: window.currentUser })) }
window.getSocket.onclose = () => { 
console.log('ok...')
document.getElementById('loginDiv').classList.remove('hidden')
document.getElementById('typingDiv').classList.add('hidden')
}
window.getSocket.onmessage = function (event) {
console.log(event.data)
if (event.data == "new game") {
 getCitation()
 return ;
}
if (event.data == "finish") {
 return ;
}
if (event.data == "top") {
 console.log('top')
 return ;
}
users = JSON.parse(event.data)
usersTable.innerHTML = ''
if (users[window.currentUser]) {
 letterInput.index = users[window.currentUser].index
}
Object.keys(users).forEach( key => {
 let prct = 100 * users[key].index / letterInput.maxIndex
 usersTable.innerHTML = usersTable.innerHTML + `<tr>
   <td style="width:1%; border: 1px solid grey">${key}</td>
   <td style="border: 1px solid grey"><div style='padding-left: ${prct}%; height: 20px; width: 3px; background-color: grey'> </div></td>
   <td style="width:1%">${users[key].index}</td>
   <td style="width:1%">${users[key].CurrentError}</td>
   </tr>\n`
})
}
}
let openTypeSocket = function() {
window.typeSocket = new WebSocket(`${socketPath}/type_ws`)
window.typeSocket.onopen = () => { window.typeSocket.send(JSON.stringify({ name: window.currentUser })) }
window.typeSocket.onclose = () => { 
console.log('asdfasdf')
document.getElementById('loginDiv').classList.remove('hidden')
document.getElementById('typingDiv').classList.add('hidden')
}
}

let start = () => {
document.getElementById('loginDiv').classList.add('hidden')
document.getElementById('typingDiv').classList.remove('hidden')
letterInput.focus()
openTypeSocket()
openGetSocket()
getCitation()
}

if (window.location.pathname.slice(1) != "") {
window.currentUser = window.location.pathname.slice(1) 
start()
}

userInput.onkeydown = e => {
if (e.which == 13) {
  userLoginButton.click()
}
}

userLoginButton.onclick = e => {
window.currentUser = userInput.value
start()
}

letterInput.oninput = e => {
let div = document.querySelector(`div.letter[index='${letterInput.index}']`)

if (div.innerHTML == e.target.value) {
div.classList.add('success')
div.classList.remove('index')
let elem = document.querySelector(`div.letter[index='${letterInput.index + 1}']`)
if (elem) {
letterInput.index += 1
elem.classList.add('index') 
}
window.typeSocket.send(JSON.stringify({
name: window.currentUser,
input: e.target.value
}))
}
else {
div.classList.add('fail')
}
e.target.value = ''
e.preventDefault()
}

// var app4 = new Vue({
//   el: '#app4',
//   data: {
//     todos: [
//       { text: 'Learn JavaScript' },
//       { text: 'Learn Vue' },
//       { text: 'Build something awesome' }
//     ]
//   }
// })

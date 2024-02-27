let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");

var socket = new WebSocket("ws://" + IP + ":8989/echo");

socket.onopen = function () {
    console.log("websocket connected");
};

socket.onmessage = function (e) {
    console.log("test");
    output.innerHTML += "Server: " + e.data + "\n";
};

socket.onerror = function (error) {
    console.log("Error in websocket");
};

function send() {
    var chat_input = document.getElementById("chat-input");
    socket.send(chat_input.value);
    chat_input.value = "";
}
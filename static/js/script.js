let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");
let socket;

function connect() {

    socket = new WebSocket("ws://" + IP + ":8989/echo");

    socket.onopen = function () {
        console.log("websocket connected at ws://" + IP + ":8989/echo");
    };

    socket.onmessage = function (event) {
        output.innerHTML += event.data;
    };

    socket.onclose = function(event) {
        console.log("WebSocket is closed.");
    };

    socket.onerror = function (error) {
        console.log("Error in websocket");
    };
}

function send() {
    var chat_input = document.getElementById("chat-input");
    socket.send(chat_input.value);
    chat_input.value = "";
}

connect();
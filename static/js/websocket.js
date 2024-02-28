let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");
let socket;
const parts = connectedUser.split(" ");
parts[0] = parts[0].substring(1);

function connect() {

    socket = new WebSocket("ws://" + IP + ":8989/echo");

    socket.onopen = function () {
        console.log("websocket connected at ws://" + IP + ":8989/echo");
    };

    socket.onmessage = function (event) {
        const messageDiv = document.createElement("div");
        messageDiv.className = "message";

        let msgData = event.data.split(" "); 

        const senderName_msg = document.createElement("div");
        senderName_msg.className = "senderName_msg";
        senderName_msg.innerHTML = msgData[0];
        messageDiv.appendChild(senderName_msg);

        const date_msg = document.createElement("div");
        date_msg.className = "date_msg";
        date_msg.innerHTML = msgData[2];
        messageDiv.appendChild(date_msg);

        const pfp_msg = document.createElement("div");
        pfp_msg.className = "pfp_msg";
        pfp_msg.innerHTML = msgData[3];
        messageDiv.appendChild(pfp_msg);

        const content_msg = document.createElement("div");
        content_msg.className = "content_msg";
        content_msg.innerHTML = msgData[4];
        messageDiv.appendChild(content_msg);
        
        output.appendChild(messageDiv);
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

    var d = new Date();
    var datestring = (d.getMonth()+1) + "/" + d.getDate() + "/" + d.getFullYear() + "_" + d.getHours() + ":" + d.getMinutes() + " ";

    socket.send(parts[0] + " receiverNameHere " + datestring + chat_input.value);
    chat_input.value = "";
}

connect();
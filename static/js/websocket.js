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
        setTimeout(() => {
        const messageDiv = document.createElement("div");
        messageDiv.className = "message";

        const senderName_msg = document.createElement("div");
        senderName_msg.className = "senderName_msg";
        senderName_msg.innerHTML = event.data.split(" ")[0];
        messageDiv.appendChild(senderName_msg);

        const pfp_msg = document.createElement("img");
        pfp_msg.className = "pfp_msg";
        pfp_msg.src = event.data.split(" ")[4];
        messageDiv.appendChild(pfp_msg);

        let temp = event.data.split(" ").slice(5).join(" ").slice(0, -1);

        const content_msg = document.createElement("div");
        content_msg.className = "content_msg";
        content_msg.innerHTML = temp;
        messageDiv.appendChild(content_msg);

        const date_msg = document.createElement("div");
        date_msg.className = "date_msg";
        date_msg.innerHTML = event.data.split(" ")[2];
        messageDiv.appendChild(date_msg);
        
        output.appendChild(messageDiv);
        output.scrollTop = output.scrollHeight;
        }, 50);
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

    socket.send(parts[0] + " testtesttest " + datestring + chat_input.value);
    chat_input.value = "";
}

connect();
setTimeout(() => {
    socket.send(parts[0]);
}, 500);

document.getElementById("chat-input").addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        console.log("test");
        event.preventDefault(); // Prevent the default form submission
        send(); // Call your send function
    }
});

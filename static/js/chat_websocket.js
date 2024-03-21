let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");
let chat_socket;
let receivername_;
let AllUser = [];
let array = [];
let AllMessages = [];
let leftovers;
let IsConnected;
let IsMsgOpen = false;
let notif = false;
let max_Messages = 10;
var coll = document.getElementsByClassName("collapsible");
const parts = connectedUser.split(" ");
parts[0] = parts[0].substring(1);

function chat_websocket() {
    chat_socket = new WebSocket("ws://" + IP + ":8989/chat_websocket");

    chat_socket.onopen = function () {
        console.log("chat_websocket connected at ws://" + IP + ":8989/chat_websocket");
    };

    chat_socket.onmessage = function (event) {
        console.log(event.data);
        setTimeout(() => {
            if (event.data === "IsCo_Yes" || event.data === "IsCo_No") {
                IsConnected = event.data;

            } else if (event.data.substring(0, 4) === "U_N ") {
                AllUser = event.data
                    .substring(4, event.data.length - 1)
                    .split(" ");

            } else {
                // Chat
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

                let temp = event.data
                    .split(" ")
                    .slice(5)
                    .join(" ")
                    .slice(0, -1);

                const content_msg = document.createElement("div");
                content_msg.className = "content_msg";
                content_msg.innerHTML = temp;
                messageDiv.appendChild(content_msg);

                const date_msg = document.createElement("div");
                date_msg.className = "date_msg";
                date_msg.innerHTML = event.data
                    .split(" ")[2]
                    .split("_")
                    .join("<br>");
                messageDiv.appendChild(date_msg);

                output.appendChild(messageDiv);
                output.scrollTop = output.scrollHeight;

                if (!IsMsgOpen && !notif) {
                    console.log("notif")
                    const notif = document.createElement("img");
                    notif.className = "notif";
                    notif.src = "../static/img/disconnected.webp";
                    collapsible = document.getElementById("collapsible");
                    collapsible.appendChild(notif);
                    notif = true;
                }
            }
        }, 50);
    };

    chat_socket.onclose = function () {
        console.log("WebSocket is closed.");
    };

    chat_socket.onerror = function (error) {
        console.log("Error in websocket");
        console.log(error);
    };
}

function send() {
    var chat_input = document.getElementById("chat-input");

    if (chat_input.value != "") {
        var d = new Date();
        var datestring =
            (d.getMonth() + 1).toString().padStart(2, "0") +
            "/" +
            d.getDate().toString().padStart(2, "0") +
            "/" +
            d.getFullYear() +
            "_" +
            d.getHours().toString().padStart(2, "0") +
            ":" +
            d.getMinutes().toString().padStart(2, "0") +
            " ";

            chat_socket.send(
            parts[0] + " " + receivername_ + " " + datestring + chat_input.value
        );
        chat_input.value = "";
    }
}

chat_websocket();
setTimeout(() => {
    chat_socket.send(parts[0]);
}, 300);

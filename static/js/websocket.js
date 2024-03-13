let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");
let socket;
let receivername_;
let AllUser = [];
let array = [];
let AllMessages = [];
let leftovers;
let IsConnected;
let max_Messages = 10;
const parts = connectedUser.split(" ");
parts[0] = parts[0].substring(1);

function connect() {
    socket = new WebSocket("ws://" + IP + ":8989/echo");

    socket.onopen = function () {
        console.log("websocket connected at ws://" + IP + ":8989/echo");
    };

    socket.onmessage = function (event) {
        setTimeout(() => {
            if (event.data === "IsCo_Yes" || event.data === "IsCo_No") {
                IsConnected = event.data;
            } else if (event.data.substring(0, 4) === "U_N ") {
                AllUser = event.data
                    .substring(4, event.data.length - 1)
                    .split(" ");
            } else {
                if (event.data.substring(0, 8) === "PUBLISH_") {
                    array = event.data.split(" ", 6);

                    // Create the post container
                    const postDiv = document.createElement("div");
                    postDiv.className = "rangepost";

                    // Create the user info container
                    const userInfoDiv = document.createElement("div");
                    userInfoDiv.className = "postuserinfo";

                    // Create the user profile picture
                    const img = document.createElement("img");
                    img.className = "pfp";
                    img.src = array[1];
                    img.alt = array[2] + `'s profile picture`;
                    userInfoDiv.appendChild(img);

                    // Create the username span
                    const usernameSpan = document.createElement("span");
                    usernameSpan.className = "userid";
                    usernameSpan.textContent = array[2];
                    userInfoDiv.appendChild(usernameSpan);

                    // Append user info to the post container
                    postDiv.appendChild(userInfoDiv);

                    // Create the category div
                    const categoryDiv = document.createElement("div");
                    categoryDiv.className = "category";
                    if (array[3] === "_&nbsp_") {
                        array[3] = "";
                    }
                    if (array[4] === "_&nbsp_") {
                        array[4] = "";
                    }
                    categoryDiv.textContent = array[3] + " " + array[4];
                    postDiv.appendChild(categoryDiv);

                    // Create the content paragraph
                    const contentP = document.createElement("p");
                    contentP.className = "content";
                    contentP.textContent = event.data
                        .split(" ")
                        .slice(5)
                        .join(" ");
                    postDiv.appendChild(contentP);

                    // Create the comment button
                    const comment_button = document.createElement("button");
                    comment_button.className = "comment_button";
                    comment_button.type = "button";
                    comment_button.onclick = () => {
                        Open_Comments(array[5]); // post id
                    };

                    const comment_img = document.createElement("img");
                    comment_img.src = "../static/img/comment_img.png";
                    comment_img.className = "comment_img";
                    comment_button.appendChild(comment_img);
                    postDiv.appendChild(comment_button);

                    postDiv.appendChild(comment_button);

                    let allPostsDiv = document.getElementById("allposts");
                    // Append the post container to the main container
                    allPostsDiv.prepend(postDiv);

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
                }
            }
        }, 50);
    };

    socket.onclose = function (event) {
        console.log("WebSocket is closed.");
    };

    socket.onerror = function (error) {
        console.log("Error in websocket");
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

        socket.send(
            parts[0] + " " + receivername_ + " " + datestring + chat_input.value
        );
        chat_input.value = "";
    }
}

function publish() {
    userId = document.getElementById("userid").getAttribute("value");
    cat1 = document.getElementById("cat1").value.replaceAll(/ /g, "_");
    cat2 = document.getElementById("cat2").value.replaceAll(/ /g, "_");
    postcontent = document.getElementById("postcontent").value;
    socket.send(
        "PUBLISH_" + " " + userId + " " + cat1 + " " + cat2 + " " + postcontent
    );

    // Clear the input fields
    document.getElementById("cat1").value = "";
    document.getElementById("cat2").value = "";
    document.getElementById("postcontent").value = "";
}

connect();
setTimeout(() => {
    socket.send(parts[0]);
}, 500);

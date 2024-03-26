let comment_socket;

function Open_Comments(postid) {
    document.querySelector("#homepage").style.display = "none";
    document.querySelector("#postpage").style.display = "block";
    document.querySelector('link[rel="stylesheet"]').href =
        "../static/CSS/postpage.css";

    comment_socket = new WebSocket("ws://" + IP + ":8989/comment_websocket");

    comment_socket.onopen = function () {
        console.log(
            "comment_socket connected at ws://" + IP + ":8989/comment_websocket"
        );
    };

    comment_socket.onmessage = function (event) {
        console.log(event.data);
        setTimeout(() => {
            if (event.data.substring(0, 4) === "P_M ") {
                array = event.data.split(" ", 6);

                // Create the post container
                const postDiv = document.createElement("div");
                postDiv.className = "rangecomment";

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

                const postId = document.createElement("div");
                postId.id = "postid";
                postId.value = array[5];
                postId.style.visibility = "hidden";
                postDiv.appendChild(postId);

                // Create the content paragraph
                const contentP = document.createElement("p");
                contentP.className = "content";
                contentP.textContent = event.data.split(" ").slice(6).join(" ");
                postDiv.appendChild(contentP);

                let commentsDiv = document.getElementById("initialpost");
                // Append the post container to the main container
                commentsDiv.append(postDiv);
            }
            document.querySelector('link[rel="stylesheet"]').href =
                "../static/CSS/postpage.css";

            if (event.data.substring(0, 4) === "C_M ") {
                array = event.data.split(" ", 6);

                // Create the post container
                const postDiv = document.createElement("div");
                postDiv.className = "rangecomment";

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

                // Create the content paragraph
                const contentP = document.createElement("p");
                contentP.className = "content";
                contentP.textContent = event.data.split(" ").slice(6).join(" ");
                postDiv.appendChild(contentP);

                let commentsDiv = document.getElementById("allcomments");
                // Append the post container to the main container
                commentsDiv.prepend(postDiv);
            }
            document.querySelector('link[rel="stylesheet"]').href =
                "/static/CSS/postpage.css";
        }, 50);
    };

    comment_socket.onclose = function () {
        console.log("comment_socket is closed.");
    };

    comment_socket.onerror = function (error) {
        console.log("Error in websocket");
        console.log(error);
    };

    setTimeout(() => {
        comment_socket.send("C_M " + parts[0] + " " + postid);
    }, 300);
}

function LoadHomePage() {
    comment_socket.close(1000, "Normal closure");
    document.querySelectorAll(".rangecomment").forEach((el) => el.remove());
    document.querySelector("#homepage").style.display = "block";
    document.querySelector("#postpage").style.display = "none";
    document.querySelector('link[rel="stylesheet"]').href =
        "/static/CSS/homepage.css";
}

function PublishComment() {
    var userId = document.getElementById("userid").getAttribute("value");
    var postcontent = document.getElementById("postcontent_").value;
    var postid = document.getElementById("postid").value;

    // Ensure the WebSocket connection is open before sending
    if (comment_socket.readyState === WebSocket.OPEN) {
        comment_socket.send(
            "P_C" + " " + postid + " " + userId + " " + postcontent
        );

        // Clear the input fields
        document.getElementById("postcontent_").value = "";
    } else {
        console.error("WebSocket connection is not open. Cannot send message.");
    }
}

let post_socket;

function post_websocket() {

    post_socket = new WebSocket("ws://" + IP + ":8989/post_websocket");

    post_socket.onopen = function () {
        console.log("post_socket connected at ws://" + IP + ":8989/post_websocket");
    };

    post_socket.onmessage = function (event) {
        setTimeout(() => {
            if (event.data.substring(0, 4) === "P_B ") {
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
                    .slice(6)
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
    
            }
        }, 50);
    };

    post_socket.onclose = function () {
        console.log("post_socket is closed.");
    };

    post_socket.onerror = function (error) {
        console.log("Error in websocket");
        console.log(error);
    };
}

function publish() {
    userId = document.getElementById("userid").getAttribute("value");
    cat1 = document.getElementById("cat1").value.replaceAll(/ /g, "_");
    cat2 = document.getElementById("cat2").value.replaceAll(/ /g, "_");
    postcontent = document.getElementById("postcontent").value;
    post_socket.send(
        "P_B" + " " + userId + " " + cat1 + " " + cat2 + " " + postcontent
    );

    // Clear the input fields
    document.getElementById("cat1").value = "";
    document.getElementById("cat2").value = "";
    document.getElementById("postcontent").value = "";
}

post_websocket();
setTimeout(() => {
    post_socket.send("1_D " + parts[0]);
}, 300);
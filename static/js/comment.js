let comment_socket;

function Open_Comments(postid) {
    document.querySelector('#homepage').style.display = 'none';
    document.querySelector('#postpage').style.display = 'block';
    document.querySelector('link[rel="stylesheet"]').href = "../static/CSS/postpage.css";

    comment_socket = new WebSocket("ws://" + IP + ":8989/comment_websocket");

    comment_socket.onopen = function () {
        console.log("comment_socket connected at ws://" + IP + ":8989/comment_websocket");
    };

    comment_socket.onmessage = function (event) {
        setTimeout(() => {
            // récupérer un comments et l'afficher
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
    document.querySelector('#homepage').style.display = 'block';
    document.querySelector('#postpage').style.display = 'none';
    document.querySelector('link[rel="stylesheet"]').href = "/static/CSS/homepage.css";
}
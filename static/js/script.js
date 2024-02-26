//var output = document.getElementById("output");
let IP = document.getElementById("IP").getAttribute("value");
let connectedUser = document
    .getElementById("connectedUser")
    .getAttribute("value");
// var socket = new WebSocket("ws://" + IP + ":8989/echo");
var cssLink = document.querySelector('link[rel="stylesheet"]');
/*socket.onopen = function () {
    //output.innerHTML += "Status: Connected to " + IP + "\n";
};

socket.onmessage = function (e) {
    //output.innerHTML += "Server: " + e.data + "\n";
};

socket.onerror = function (error) {
    //output.innerHTML += "Error connecting to " + IP + "\n";
};*/

if (connectedUser != "") {
    console.log("oktamer");
    document.getElementById("forms-container").innerHTML = "";
}

//var output = document.getElementById("output");
var socket = new WebSocket("ws://" + IP + ":8989/echo");
var cssLink = document.querySelector('link[rel="stylesheet"]'); // do not use this in the forum
socket.onopen = function () {
    //output.innerHTML += "Status: Connected to " + IP + "\n";
};

socket.onmessage = function (e) {
    //output.innerHTML += "Server: " + e.data + "\n";
};

socket.onerror = function (error) {
    //output.innerHTML += "Error connecting to " + IP + "\n";
};


document.getElementById("chat-input").addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        event.preventDefault(); // Prevent the default form submission
        send(); // Call your send function
    }
});

function render_chat(nickname) {
    receivername_ = nickname
    
    let header_chat = document.getElementById("header_chat");

    let conversations = document.getElementsByClassName("conversation_with");
    while (conversations.length > 0) {
        conversations[0].parentNode.removeChild(conversations[0]);
    }    

    const conversation = document.createElement("span");
    conversation.className = "conversation_with";
    conversation.innerHTML = "Conversation with " + nickname
    header_chat.appendChild(conversation);

    document.querySelector('#output').innerHTML = "";
    // ask the websocket the conversation
    socket.send("GAM " + nickname + " " + parts[0] + " _");

    document.querySelector('.msg-inputs').style.display = 'block';
    document.querySelector('#output').style.display = 'block';
    document.querySelector('.header_chat').style.display = 'flex';
    document.querySelector('.nickname_button_div').style.display = 'none';
}

function hide_chat() {
    document.querySelector('.msg-inputs').style.display = 'none';
    document.querySelector('#output').style.display = 'none';
    document.querySelector('.header_chat').style.display = 'none';
    document.getElementById("nickname_button_div").innerHTML = "";

    // ask the websocket all the user
    socket.send("U_N " + parts[0]);
    setTimeout(() => {
        AllUser.forEach(name => {
            if (name != parts[0]) {
                const nickname_button = document.createElement("button");
                nickname_button.className = "nickname_button";
                nickname_button.type = "button";
                nickname_button.onclick = () => {
                    render_chat(name);
                };
                nickname_button.innerHTML = name;
                document.getElementById("nickname_button_div").appendChild(nickname_button);
            };
        });
    }, 100);

    document.querySelector('.nickname_button_div').style.display = 'block';
}

setTimeout(() => {
    hide_chat()
}, 500);
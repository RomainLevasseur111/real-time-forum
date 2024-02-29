
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
    socket.send("GAM " + nickname + " " + parts[0] + " _");

    document.querySelector('.msg-inputs').style.display = 'block';
    document.querySelector('#output').style.display = 'block';
    document.querySelector('.header_chat').style.display = 'flex';
    var elements = document.querySelectorAll('.nickname_button');
    elements.forEach(function(element) {
        element.style.display = 'none';
    });
}

function hide_chat() {
    document.querySelector('.msg-inputs').style.display = 'none';
    document.querySelector('#output').style.display = 'none';
    document.querySelector('.header_chat').style.display = 'none';
    var elements = document.querySelectorAll('.nickname_button');
    elements.forEach(function(element) {
        element.style.display = 'block';
    });
}

hide_chat();
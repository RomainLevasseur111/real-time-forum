
document.getElementById("chat-input").addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        event.preventDefault(); // Prevent the default form submission
        send(); // Call your send function
    }
});

function render_chat(nickname) {
    receivername_ = nickname
    socket.send("IsCo " + nickname + " " + parts[0]);
    setTimeout(() => {
        if (IsConnected === "IsCo_Yes") {
            IsConnected = "";
        
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
    
            setTimeout(() => {
                hide_messages();
            }, 75);
    
            document.querySelector('.msg-inputs').style.display = 'block';
            document.querySelector('#output').style.display = 'block';
            document.querySelector('.header_chat').style.display = 'flex';
            document.querySelector('.nickname_button_div').style.display = 'none';
        }
    }, 100);
}

function hide_chat() {
    document.querySelector('.msg-inputs').style.display = 'none';
    document.querySelector('#output').style.display = 'none';
    document.querySelector('.header_chat').style.display = 'none';
    max_Messages = 10;

    // ask the websocket all the user
    socket.send("U_N " + parts[0]);

    setTimeout(() => {
        document.getElementById("nickname_button_div").innerHTML = "";
        AllUser.forEach((user, i) => {
            if (user != parts[0] && (i)%3 === 0) {
                const nickname_button = document.createElement("button");
                nickname_button.className = "nickname_button";
                nickname_button.type = "button";
                nickname_button.onclick = () => {
                    render_chat(user);
                };
                const conn_img = document.createElement("img");
                conn_img.src = AllUser[i+2];
                if (AllUser[i+2] === "../static/img/connected.png") {
                    conn_img.style.height = "17px";
                    conn_img.style.width = "17px";
                };
                conn_img.className = "conn_img";
                const pfp_chat_img = document.createElement("img");
                pfp_chat_img.src = AllUser[i+1];
                pfp_chat_img.className = "pfp_chat_img"
                nickname_button.appendChild(pfp_chat_img);
                nickname_button.innerHTML += user;
                nickname_button.appendChild(conn_img);
                document.getElementById("nickname_button_div").appendChild(nickname_button);
            };
        });
    }, 100);

    document.querySelector('.nickname_button_div').style.display = 'grid';
    setTimeout(() => {
        if (document.querySelector('.msg-inputs').style.display === 'none') {
            hide_chat();
        }
    }, 2500);
}


function hide_messages() {
    let height = 0;

    let messageCount = output.childElementCount;
    for (let i = 0; i < messageCount; i++) {
        if (i >= messageCount - max_Messages) {
            if (output.children[i].style.display != 'flex') {
                output.children[i].style.display = 'flex';
                height += output.children[i].offsetHeight + 10;
            }
        } else {
            output.children[i].style.display = 'none';
        }
    }
    output.scrollTop = height;
}

function debounce(func, wait) {
    let timeout;
    return function() {
        const context = this, args = arguments;
        clearTimeout(timeout);
        timeout = setTimeout(function() {
            func.apply(context, args);
        }, wait);
    };
}

const output = document.getElementById('output');

const handleScroll = debounce(function() {
    if (output.scrollTop === 0) {
        max_Messages += 10;
        hide_messages();
    }
}, 500);

output.addEventListener('scroll', handleScroll);

setTimeout(() => {
    hide_chat()
}, 500);

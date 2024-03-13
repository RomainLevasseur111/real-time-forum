function Open_Comments(postid) {
    document.querySelector('#homepage').style.display = 'none';
    document.querySelector('#postpage').style.display = 'block';
    document.querySelector('link[rel="stylesheet"]').href = "../static/CSS/postpage.css";
}
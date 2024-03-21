var i;

for (i = 0; i < coll.length; i++) {
  coll[i].addEventListener("click", function() {
    this.classList.toggle("collapsible-active");
    var content = this.nextElementSibling;
    if (content.style.width){
      content.style.width = null;
      IsMsgOpen = false;
      notif = false;
    } else {
      content.style.width = "29vw";
      IsMsgOpen = true;
    }
  });
}
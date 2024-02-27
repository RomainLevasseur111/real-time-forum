var coll = document.getElementsByClassName("collapsible");
var i;

for (i = 0; i < coll.length; i++) {
  coll[i].addEventListener("click", function() {
    this.classList.toggle("collapsible-active");
    var content = this.nextElementSibling;
    if (content.style.width){
      content.style.height = null;
      content.style.width = null;
      content.style.minWidth = null;
    } else {
      content.style.width = "50vw";
      setTimeout(() => {
        if (content.style.width === "50vw") {
            content.style.minWidth = "300px";
        }
    },  2000);
    }
  });
}
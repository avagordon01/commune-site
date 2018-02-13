function set_freshness(value) {
    var d = new Date();
    d.setTime(d.getTime() + (12*60*60*1000));
    document.cookie = "freshness=" + value + ";expires=" + d.toUTCString() + ";path=/;Secure";
    location.reload(true)
}

function toggle_hidden(element) {
    if (element.innerHTML == "hide") {
        element.innerHTML = "show"
        els = element.parentNode.querySelectorAll("*")
        for (var i = 0; i < els.length; i++) {
            if (els[i] != element) {
                els[i].style.display = "none"
            }
        }
    } else {
        els = element.parentNode.querySelectorAll("*")
        for (var i = 0;i < els.length; i++) {
            if (els[i] != element) {
                els[i].style.display = null
            }
        }
        element.innerHTML = "hide"
    }
}

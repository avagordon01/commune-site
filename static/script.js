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

function toggle_comment(element) {
    var comment = document.querySelector("#comment")
    if (comment == element.nextSibling) {
        comment.style.display = null
        post = document.querySelector(".post")
        post.parentElement.insertBefore(comment, post.nextSibling)
    } else {
        comment.style.display = "initial"
        post = document.querySelector(".post")
        comment.post_id.value = post.id
        if (element.parentElement != post) {
            comment.comment_id.value = element.parentElement.id
        } else {
            comment.comment_id.value = ""
        }
        element.parentElement.insertBefore(comment, element.nextSibling)
        comment.text.focus()
    }
}

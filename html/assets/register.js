const xhr = new XMLHttpRequest();
var result = "";

function submit() {
    let user = document.getElementById("usernameField").value;
    let params = {"username": user};
    xhr.open("POST", "/register");
    xhr.send(JSON.stringify(params));
}

xhr.onload = () => {
    if (xhr.status == 200) {
        result = xhr.response;
        result = result.replace("\n", "<br>");
        var msg = JSON.parse(result);

        document.getElementById("resultField").innerHTML = "<h3 style=\"color: #EEE;\">Username " + msg.msg + "</h3>";
    } else {
        console.error("Error!");
    }
};
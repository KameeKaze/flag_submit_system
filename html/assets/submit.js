const xhr = new XMLHttpRequest();
var result = "";

function submit() {
    let token = document.getElementById("tokenField").value;
    let flag = document.getElementById("flagField").value;
    let params = {"token": token, "flag": flag};
    xhr.open("POST", "/submit");
    xhr.send(JSON.stringify(params));
}

xhr.onload = () => {
    if (xhr.status == 200) {
        result = xhr.response;
        result = result.replace("\n", "<br>");
        var msg = JSON.parse(result);

        document.getElementById("resultField").innerHTML = "<h3 style=\"color: #EEE;\">" + msg.msg + "</h3>";
    } else {
        console.error("Error!");
    }
};
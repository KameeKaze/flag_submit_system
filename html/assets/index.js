const xhr = new XMLHttpRequest();
var result = "";

function getScoreboard() {
    xhr.open("POST", "/scoreboard");
    xhr.send();
}

xhr.onload = () => {
    if (xhr.status == 200) {
        result = xhr.response;
        result = result.replace("\n", "<br>");
        console.log(msg);
        var msg = JSON.parse(result);

        for (var i = 0; i < msg.length; i++) {
            document.getElementById("scoreboard").innerHTML = document.getElementById("scoreboard").innerHTML+
            "<div class=\"row\">"+
                "<div class=\"name\">"+
                "<h5>"+
                String(i+1) + "&nbsp;&nbsp;&nbsp;"+
                msg[i].Username+"</h4>"+
                "</div><div class=\"score\">"+
                "<h5>"+msg[i].score+"</h4>"+
                "</div></div>";
        }
    } else {
        console.error("Error!");
    }
};
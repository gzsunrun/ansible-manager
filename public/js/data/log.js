var term = new Terminal({
    cols: 80,
    rows: 40,
    screenKeys: false,
    useStyle: true,
    cursorBlink: true,
    convertEol: true
  });
// var sdata
// term.on('data', function(data){
//     if (event.keyCode == 13) {
//         alert(sdata);
//         sdata =""
//         term.write("\n")
//     }
//     term.write(data)
//     sdata +=data
// })
term.open($("#term").empty()[0]);


function getQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]); return null;
}
var task_id=getQueryString("task_id")
window.onload = function () {
var conn;
var msg = document.getElementById("msg");
var log = document.getElementById("log");

function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}


if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/api/ansible/ws?task_id="+task_id);
    conn.onclose = function (evt) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        appendLog(item);
    };
    conn.onmessage = function (evt) {
        var obj = JSON.parse(evt.data)
        if (obj.type=="log"){
            data=obj.output
            if(data.indexOf("failed")==0||data.indexOf("fatal")==0)
            {
                term.write("\033[31m"+data+'\033[0m\r\n');
                return
            }
            if(data.indexOf("changed")==0||data.indexOf("warnning")==0)
            {
                term.write("\033[33m"+data+'\033[0m\r\n');
                return
            }
            if(data.indexOf("ok")==0)
            {
                term.write("\033[32m"+data+'\033[0m\r\n');
                return
            }
            term.write(data+'\r\n');
            
        }
        if (obj.type=="update"){
            term.write(`start:`+obj.start+'\r\n');
            term.write(`end:`+obj.end+'\r\n');
            term.write(`status:`+obj.status+'\r\n');
        }
    };
} else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
}
};
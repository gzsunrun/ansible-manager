
$(function(){
    FindTimers()
    setInterval("FindTimers()",5000)
})

function FindTimers(){
    var success=function(msg){
        $("#timer-list").empty()
        $.each(msg, function (i, val) {
            if (val.timer_status){
                var status=`<a href="javascript:StopTimer('`+val.timer_id+`');" class="btn btn-xs btn-warning" id="stop-`+val.timer_id+`" title="停止">
                    <i class="icon icon-stop"></i>停止
                    </a>`
            }else{
                var status=`<a href="javascript:StartTimer('`+val.timer_id+`');" class="btn btn-xs btn-success" id="start-`+val.timer_id+`" title="启动">
                    <i class="icon icon-play"></i>启动
                    </a>`
            }
            var h="--"
            var m="--"
            var s="--"
            if (val.timer_status){
                h=parseInt(val.timer_surplus/3600)
                m=parseInt((val.timer_surplus-h*3600)/60)
                s=val.timer_surplus-h*3600-m*60
            }
            
            
            var html = `<tr>
            <td>`+ val.timer_id + `</td>
            <td>`+ val.timer_name + `</td>
            <td>`+ val.timer_interval + `</td>
            <td>`+ h +"小时"+m+"分钟"+s+"秒"+ `</td>
            <td>`+ val.timer_repeat + `</td>
            <td>`+ val.timer_status + `</td>
            <td>`+  new Date(val.created).format("yyyy-MM-dd hh:mm:ss") + `</td>
            <td>
                `+status+`
                <a href="javascript:DelTimer('`+val.timer_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#timer-list').append(html)
        })
    };
    AjaxReq(
        "get",
        "../ansible/common/timer/list",
        {},
        function () { },
        success,
        ReqErr
    );
}

function StartTimer(id){
    AjaxReq(
        "get",
        "../ansible/common/timer/start",
        {timer_id:id},
        function () { },
        ReqSuccess,
        ReqErr
    );
}

function StopTimer(id){
    AjaxReq(
        "get",
        "../ansible/common/timer/stop",
        {timer_id:id},
        function () { },
        ReqSuccess,
        ReqErr
    );
}

function DelTimer(id){
    AjaxReq(
        "get",
        "../ansible/common/timer/del",
        {timer_id:id},
        function () { },
        function () {
            ReqSuccess
            FindTimers()
         },
        ReqErr
    );
}
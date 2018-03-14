

$(function(){
    AjaxReq(
        "get",
        "../ansible/common/task/count",
        {},
        function () { },
        function(msg){
            $("#total").html(msg.all_total)
            $("#run-total").html(msg.run_total)
            $("#success-total").html(msg.success_total)
            $("#error-total").html(msg.error_total)
        },
        ReqErr
    );
    GetNodeList()
})


function GetNodeList() {
    var success = function (msg) {
        $("#node-list").empty()
        $.each(msg, function (i, val) {
            var html = `<tr>
            <td>`+ val.node_id + `</td>
            <td>`+ val.node_ip + `</td>
            <td>`+ val.node_port + `</td>
            <td>`+ val.node_master + `</td>
            <td>`+ val.node_worker + `</td>
        </tr>`
            $('#node-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/nodes",
        {},
        function () { },
        success,
        ReqErr
    );
}
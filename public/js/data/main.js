

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
})
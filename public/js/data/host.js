
$(function(){
    GetHostList()
})

$("#host-commit").click(function(){
    EditHost()
})

$("#host-add").click(function(){
    CleanHostForm()
})

// get host list
function GetHostList() {
    $("#host-list").html(`<tr><td colspan=6><i class="icon icon-spin icon-spinner-indicator icon-3x"></i></td></tr>`)
    var success = function (msg) {
        $("#host-list").empty()
        $.each(msg, function (i, val) {
            if (val.host_status) {
                var status = `<span style="color:green" class="icon icon-check-circle-o"></span>success`
            } else {
                var status = `<span style="color:red" class="icon icon-remove-circle"></span>fail`
            }
            var html = `<tr>
            <td>`+ val.host_id + `</td>
            <td>`+ val.host_alias + `</td>
            <td>`+ val.host_ip + `</td>
            <td>`+ status + `</td>
            <td>`+ val.created + `</td>
            <td>
                <a href="javascript:GetHost('`+val.host_id+`');" class="btn btn-xs btn-primary">编辑</a>
                <a href="javascript:DelHost('`+val.host_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#host-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/hosts",
        {},
        function () { },
        success,
        ReqErr
    );
}

// update host
function EditHost() {
    var data = {
        host_id: $('#host-id').val(),
        host_alias: $('#host-alias').val(),
        host_name: $('#host-name').val(),
        host_ip: $('#host-ip').val(),
        host_user: $('#host-user').val(),
        host_password: $('#host-password').val(),
        host_key: $('#host-key').val()
    }
    AjaxReq(
        "post",
        "../ansible/common/hosts/create",
        JSON.stringify(data),
        function () { },
        function () {
            GetHostList()
            ReqSuccess()
            CleanHostForm()
            $('#host-modal').modal('hide');
        },
        ReqErr
    )
}

// delete host
function DelHost(id) {
    AjaxReq(
        "get",
        "../ansible/common/hosts/del",
        {host_id:id},
        function () { },
        function () {
            GetHostList()
            ReqSuccess()
        },
        ReqErr
    )
    
}

function GetHost(id){
    AjaxReq(
        "get",
        "../ansible/common/hosts/get",
        {host_id:id},
        function () { },
        function (msg) {
            $('#host-id').val(msg.host_id)
            $('#host-alias').val(msg.host_alias)
            $('#host-name').val(msg.host_name)
            $('#host-ip').val(msg.host_ip)
            $('#host-user').val(msg.host_user)
            $('#host-modal').modal('show');
        },
        ReqErr
    )
}

function CleanHostForm(){
    $('#host-id').val("")
    $('#host-alias').val("")
    $('#host-name').val("")
    $('#host-ip').val("")
    $('#host-user').val("")
}
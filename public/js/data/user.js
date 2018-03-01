$(function(){
    GetUserList()
})

$("#user-commit").click(function(){
    EditUser()
})

$("#user-add").click(function(){
    CleanUserForm()
})

function GetUserList() {
    var success = function (msg) {
        $("#user-list").empty()
        $.each(msg, function (i, val) {
            var html = `<tr>
            <td>`+ val.user_id + `</td>
            <td>`+ val.user_account + `</td>
            <td>`+ val.created + `</td>
            <td>
                <a href="javascript:GetUser('`+val.user_id+`','`+val.user_account+`');" class="btn btn-xs btn-primary">编辑</a>
                <a href="javascript:DelUser('`+val.user_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#user-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/user",
        {},
        function () { },
        success,
        ReqErr
    );
}


function EditUser() {
    var data = {
        user_id: $('#user-id').val(),
        user_account: $('#user-account').val(),
        user_password: $('#user-password').val(),
    }
    AjaxReq(
        "post",
        "../ansible/common/user/create",
        JSON.stringify(data),
        function () { },
        function () {
            GetUserList()
            ReqSuccess()
            CleanUserForm()
            $('#user-modal').modal('hide');
        },
        ReqErr
    )
}

function GetUser(id,account){
    $('#user-id').val(id)
    $('#user-account').val(account)
    $('#user-account').attr("disabled",true)
    $('#user-modal').modal('show');
}

// delete user
function DelUser(id) {
    AjaxReq(
        "get",
        "../ansible/common/user/del",
        {uid:id},
        function () { },
        function () {
            GetUserList()
            ReqSuccess()
        },
        ReqErr
    )
    
}


function CleanUserForm(){
    $('#user-id').val("")
    $('#user-account').val("")
    $('#user-password').val("")
    $('#user-account').attr("disabled",false)
}
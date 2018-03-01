
$(function(){
    GetRepoList()
})

$("#upload").click(function(){
    CreateRepo()
})

$("#repo-add").click(function(){
    CleanRepoForm()
})

// get host list
function GetRepoList() {
    var success = function (msg) {
        $("#repo-list").empty()
        $.each(msg, function (i, val) {
            var html = `<tr>
            <td>`+ val.repo_id + `</td>
            <td>`+ val.repo_name + `</td>
            <td>`+ val.created + `</td>
            <td>
                <a href="javascript:DelRepo('`+val.repo_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#repo-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/repo",
        {},
        function () { },
        success,
        ReqErr
    );
}

// update host
function CreateRepo(){
    var fd = new FormData();
    fd.append("repo_path",$('#repo-path')[0].files[0]);
    fd.append("repo_name",$("#repo-name").val());
    fd.append("repo_desc",$("#repo-desc").val());
   if ($("#repo-name").val()==""||$('#repo-path')[0].files[0]==null){
       alert("脚本名不能为空")
       return
   }
    var xhr = new XMLHttpRequest();
    if ( xhr.upload ) {
        $("#upload").addClass("disabled")
    }
    xhr.onreadystatechange = function(e) {
        if ( 4 == this.readyState ) {
            if (xhr.status == 204) {
                $("#repo-name").val("")
                $("#repo-path").val("")
                $("#repo-parse").val("")
                $('#repo-modal').modal('hide');
                GetRepoList()    
            }else{
                new $.zui.Messager("上传失败", {
                        type: 'danger' 
                        }).show();
                $("#upload").removeClass("disabled")
            }
        }
    };
    xhr.open('post', '../ansible/common/repo/create', true);
    xhr.setRequestHeader('Authorization', $.cookie("Auth"));
    xhr.send(fd);
}

// delete host
function DelRepo(id) {
    AjaxReq(
        "get",
        "../ansible/common/repo/delete",
        {repo_id:id},
        function () { },
        function () {
            GetRepoList()
            ReqSuccess()
        },
        ReqErr
    )
    
}

function CleanRepoForm(){
    $('#repo-name').val("")
    $('#repo-path').val("")
    $('#repo-desc').val("")
}
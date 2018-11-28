var GITSTATUS

$(function(){
    GitStatus()
    GetRepoList()
})

$("#upload").click(function(){
    if (GITSTATUS){
        AddGit()
    }else{
        CreateRepo()
    }
    
})

$("#repo-add").click(function(){
    CleanRepoForm()
})

// get host list
function GetRepoList() {
    var success = function (msg) {
        $("#repo-list").empty()
        $.each(msg, function (i, val) {

        var html =`<div class="col-md-4 col-sm-6 col-lg-3">
        <span class="card">
            <a href="###">
                <p style="margin: 0 auto; width: 157px; height: 157px">
                    <img width="100%" height="100%"onerror='this.src="img/default.png"' src="../ansible/static/icon?id=`+val.repo_path+`" title="点击部署">
                </p>
            </a>
            <div class="card-heading">
                    <strong>`+ val.repo_name + `</strong>
                    <select  id="exampleInputAddress7">
                        <option>`+ val.repo_version + `</option>
                    </select>
            </div>
            <div class="card-content text-muted" style="height:60px;">`+ val.repo_desc + `</div>
            <div class="card-actions">
               <button type="button" class="btn"><i class="icon-play-circle"></i> 部署</button>
                <div class="pull-right text-danger">
                        <a class="text-danger" href="javascript:DelRepo('`+val.repo_id+`');"><i class="icon icon-trash"></i> 删除</a>
                </div>
            </div>
        </span>
    </div>`
            $('#repo-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/repo",
        {repo_type:$("#repo-type").val()},
        function () { },
        success,
        ReqErr
    );
}

function GitStatus(){
    AjaxReq(
        "get",
        "../ansible/common/repo/git/status",
        {},
        function () { },
        function (msg) { 
            GITSTATUS=msg.status
            $("#repo-add").removeClass("disabled")
            if (msg.status){
                $("#repo-url").removeClass("hide")
            }else{
                $("#repo-file").removeClass("hide")
            }
        },
        ReqErr
    );
}

// update host
function CreateRepo(){
    var fd = new FormData();
    fd.append("repo_path",$('#repo-path')[0].files[0]);
    fd.append("repo_name",$("#repo-name").val());
    fd.append("repo_type",$("#repo-type").val());
    fd.append("repo_desc",$("#repo-desc").val());
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
                $("#upload").removeClass("disabled")
                GetRepoList()    
            }else{
                new $.zui.Messager("上传失败:"+xhr.responseText, {
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

function AddGit(){
    $("#upload").addClass("disabled")
    var git_url=$("#repo-git").val()
    AjaxReq(
        "get",
        "../ansible/common/repo/git/sync",
        {git_url:git_url},
        function () { },
        function () {
            GetRepoList()
            $('#repo-modal').modal('hide');
            $("#upload").removeClass("disabled")
            ReqSuccess()
        },
        function () { 
            ReqErr()
            $("#upload").removeClass("disabled")
        }
    )
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
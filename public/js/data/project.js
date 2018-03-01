
$(function(){
    GetProjectList()
    GetHostList()
})

$("#project-commit").click(function(){
    EditProject()
})

$("#project-add").click(function(){
    CleanProjectForm()
})

// get project list
function GetProjectList() {
    var success = function (msg) {
        $("#project-list").empty()
        $.each(msg, function (i, val) {
            var html = `<tr>
            <td>`+ val.project_id + `</td>
            <td>`+ val.project_name + `</td>
            <td>`+ val.created + `</td>
            <td>
                <a href="task.html?project_id=`+val.project_id+`" class="btn btn-xs btn-success">创建任务</a>
                <a href="javascript:GetProject('`+val.project_id+`');" class="btn btn-xs btn-primary">编辑</a>
                <a href="javascript:DelProject('`+val.project_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#project-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/project",
        {},
        function () { },
        success,
        ReqErr
    );
}

// get host list
function GetHostList() {
    var success = function (msg) {
        $("#host-list").empty()
        $.each(msg, function (i, val) {
            var html = `<div class="checkbox-custom checkbox-primary">
            <input type="checkbox" name="chk" id="`+val.host_id+`" value="`+val.host_id+`">
            <label for="inputUnchecked">`+val.host_alias+`（`+val.host_ip+`）</label>
        </div>`
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
function EditProject() {
    var project_hosts=new Array;
    $('input[type="checkbox"][name="chk"]:checked').each(function() {
        project_host={
            project_id:$('#project-id').val(),
            host_id:$(this).val()
        }
        project_hosts.push(project_host);   
    });
    var project={
        project_id: $('#project-id').val(),
        project_name: $('#project-name').val()
    }
    var data = {
        project: project,
        project_hosts:project_hosts
    }
    AjaxReq(
        "post",
        "../ansible/common/project/caa",
        JSON.stringify(data),
        function () { },
        function () {
            GetProjectList()
            ReqSuccess()
            CleanProjectForm()
            $('#project-modal').modal('hide');
        },
        ReqErr
    )
}

// delete host
function DelProject(id) {
    AjaxReq(
        "get",
        "../ansible/common/project/del",
        {project_id:id},
        function () { },
        function () {
            GetProjectList()
            ReqSuccess()
        },
        ReqErr
    )
    
}

function GetProject(id){
    AjaxReq(
        "get",
        "../ansible/common/project/get",
        {project_id:id},
        function () { },
        function (msg) {
            $('#project-id').val(msg.project.project_id)
            $('#project-name').val(msg.project.project_name)
            if (msg.project_hosts!=null){
                $.each(msg.project_hosts, function (i, val) {
                    $('#'+val.host_id).attr("checked","")
                })
            }
            $('#project-modal').modal('show');
        },
        ReqErr
    )
}

function CleanProjectForm(){
    $('#project-id').val("")
    $('#project-name').val("")
    GetHostList()
}
var project_id=getQueryString("project_id")
$(function(){
    GetProjectList()
})
// get release list
function GetHelmList() {
    var project_id=$("#project-list").val()
    $("#release-list").empty()
    var success = function (msg) {
        $.each(msg, function (i, val) {
            var html = `<tr>
            <td>`+ val.name + `</td>
            <td>`+ val.app_version + `</td>
            <td>`+ val.revison + `</td>
            <td>`+ val.chart + `</td>
            <td>`+ val.namespace + `</td>
            <td>`+ val.status + `</td>
            <td>`+ val.upadted + `</td>
            <td>
                <a href="javascript:StatusView('`+val.name+`','`+project_id+`');" class="btn btn-xs btn-success">详情</a>
                <a href="javascript:UpgradeView('`+val.name+`','`+val.chart +`','`+val.namespace+`');" class="btn btn-xs btn-primary">升级</a>
                <a href="javascript:DeleteView('`+val.name+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#release-list').append(html)
        })
    };


    AjaxReq(
        "get",
        "../ansible/common/helm",
        {project_id:project_id},
        function () { },
        success,
        function(){
            console.log("aaaa")
        }
    );
}

// get project list
function GetProjectList() {
    var success = function (msg) {
        $("#project-list").empty()
        $.each(msg, function (i, val) {
            if (val.project_name.indexOf("@k8s")!=-1){
                var html = `<option value="`+val.project_id+`">`+val.project_name+`</option>`
                $('#project-list').append(html)
            }
        })
        if (project_id!=null){
            $("#project-list").val(project_id)
        }
        
        GetHelmList()
         $("#project-list").change(function(){
               GetHelmList()
         })
    };

    var err =function(){
        new $.zui.Messager('提示消息：请求失败', {
            type: 'danger' 
        }).show();
    }
    AjaxReq(
        "get",
        "../ansible/common/project",
        {},
        function () { },
        success,
        err
    );
}

function DeleteView(name){
    $("#delete-name").text(name)
    $("#delete-modal").modal("show")
}

$("#delete-release").click(function(){
    var project_id=$("#project-list").val()
    var name=$("#delete-name").text()
    var purge=new Boolean($("#delete-purge").val())
    AjaxReq(
        "delete",
        "../ansible/common/helm",
        JSON.stringify({project_id:project_id,name:name,purge:purge}),
        function () { },
        function(){
            ReqSuccess()
            GetHelmList()
            $("#delete-name").text("")
            $("#delete-modal").modal("hide")
        },
        function(XMLHttpRequest, textStatus, errorThrown){
            new $.zui.Messager("删除失败:"+textStatus, {
                type: 'danger',
                time: 10000
                }).show();
        }
    );
})

function UpgradeView(name,chart,namespace){
    AjaxReq(
        "get",
        "../ansible/common/harbor/charts",
        {},
        function () { },
        function(msg){
            $.each(msg, function (i, val) {
                var versionHtlm=""
                var find=false
                $.each(val.charts,function(i,v){
                    versionHtlm  +=`<option value="`+v.version+`">`+ v.appVersion+"/"+v.version + `</option>`
                    if (chart==val.name+"-"+v.version){
                        find=true
                    }
                })
                if (find){
                    $("#chart-version").append(versionHtlm)
                    $("#release-name").val(name)
                    $("#release-namespace").val(namespace)
                    $("#chart-name").val(val.name)
                    $("#upgrade-modal").modal("show")
                    return false
                }
            }) 
        },
        function(XMLHttpRequest, textStatus, errorThrown){
            new $.zui.Messager("获取chart列表失败:"+textStatus, {
                type: 'danger',
                time: 10000
                }).show();
        }
    );
}

$("#upgrade-release").click(function(){
    var name=$("#release-name").val()
    var namespace=$("#release-namespace").val()
    var project_id =$("#project-list").val()
    var chart_name =$("#chart-name").val()
    var chart_version=$("#chart-version").val()
    window.location.href="helm-install.html?chart="+chart_name+"&version="+chart_version+"&name="+name+"&project="+project_id+"&namespace="+namespace+"&upgrade=true";
})

function StatusView(name,project_id){
    $("#status-name").text(name)
    $("#pod-list").empty()
    var success = function (msg) {
        $.each(msg.pods_status, function (i, val) {
            var html = `<tr>
            <td>`+ val.name + `</td>
            <td>`+ val.ready + `</td>
            <td>`+ val.status + `</td>
            <td>`+ val.restart + `</td>
            <td>`+ val.age + `</td>
        </tr>`
            $('#pod-list').append(html)
        })
        $("#status-modal").modal("show")
    };

    AjaxReq(
        "get",
        "../ansible/common/helm/status",
        {project_id:project_id,name:name},
        function () { },
        success,
        function(){
            console.log("aaaa")
        }
    );
}

function Msg(){
    alert("等待完善")
}


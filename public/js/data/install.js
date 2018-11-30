
var chart=getQueryString("chart")
var version=getQueryString("version")
var release_name=getQueryString("name")
var release_namespace=getQueryString("namespace")
var release_project=getQueryString("project")
var upgrade=getQueryString("upgrade")
$(function(){
    $("#chart-name").val(chart)
    $("#chart-version").val(version)
    if (release_name!=null){
        $("#release-name").val(release_name)
        $("#release-name").attr("disabled","true")
    }
    if (release_namespace!=null){
        $("#namespace").val(release_namespace)
        $("#namespace").attr("disabled","true")
    }
    if (release_project!=null){
        $("#project-list").attr("disabled","true")
    }
    GetProjectList()
    if (upgrade){
        GetHistoryValues(release_name,release_project)
    }else{
        GetValues()
    }
    
})
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
        if (release_project!=null){
            $("#project-list").val(release_project)
        }
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
 function GetHostList(project_id){
    AjaxReq(
        "get",
        "../ansible/common/project/hosts",
        {project_id:project_id},
        function () { },
        function(msg){
            HOSTS=msg
        },
        ReqErr
    );
}

function GetValues(){
    AjaxReq(
        "get",
        "../ansible/common/helm/values",
        {chart_name:chart,chart_version:version},
        function () { },
        function(msg){
            console.log(msg)
            editor.setValue(msg.files["values.yaml"])
        },
        ReqErr
    );
}

function GetHistoryValues(name,project_id){
    AjaxReq(
        "get",
        "../ansible/common/helm/hvalues",
        {name:name,project_id:project_id},
        function () { },
        function(msg){
            console.log(msg)
            editor.setValue(msg["values"])
        },
        ReqErr
    );
}

function Install(){
    var name=$("#release-name").val()
    var project_id =$("#project-list").val()
    var chart_name =$("#chart-name").val()
    var chart_version=$("#chart-version").val()
    var namespace=$("#namespace").val()
    var update =new Boolean(upgrade)
    var values=editor.getValue()
    $("#release-btn").addClass("disabled")
    AjaxReq(
        "post",
        "../ansible/common/helm",
        JSON.stringify({chart_name:chart_name,chart_version:chart_version,name:name,project_id:project_id,namespace:namespace,values:values,update:update}),
        function () { },
        function(msg){
            window.location.href="helm-status.html?project_id="+project_id
        },
        function(XMLHttpRequest, textStatus, errorThrown){
            new $.zui.Messager("发布失败:"+XMLHttpRequest.responseText , {
                type: 'danger',
                time: 10000
            }).show();
            $("#release-btn").removeClass("disabled")
        }
    );

}

$("#release-btn").click(function(){
    Install()
})
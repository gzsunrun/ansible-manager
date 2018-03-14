var project_id=getQueryString("project_id")

var HOSTS 
var REPOVARS
var TASKVARS={
    task_vars:[],
    task_group:[],
    task_id:"",
    task_name:"",
    task_tag:"",
    project_id:"",
    repo_id:""
}
$(function(){
    TaskList()
    GetHostList()
    GetRepoList()
})

$("#repo-id").change(function(){
	var repo_id=$('#repo-id').val()
	if (repo_id!=""){
		$('#task-create').removeClass("disabled")
		GetRepoVars(repo_id,false)
	}else{
		$('#task-create').addClass("disabled")
		CleanTask()
	}
})

$("#task-add").click(function(){
	$('#task-create').addClass("disabled")
	CleanTask()
})

 // get host list
function GetHostList(){
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

function GetRepoList() {
    var success = function (msg) {
        $("#repo-list").empty()
        $('#repo-list').append(`<option value="">请选择脚本</option>`)
        $.each(msg, function (i, val) {
            var html = `<option value="`+val.repo_id+`">`+ val.repo_name + `</option>`
            $('#repo-id').append(html)
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

function GetRepoVars(id,update){
    AjaxReq(
        "get",
        "../ansible/common/vars",
        {repo_id:id},
        function () { },
        function(msg){
			REPOVARS=msg
            TASKVARS.task_vars=msg.vars
            createInvView()
			createVarsView()
			if (update){
				editInvView()
            	$('#task-modal').modal('show');
			}
        },
        ReqErr
    );
}

function GetTask(id){
    AjaxReq(
        "get",
        "../ansible/common/task/get",
        {task_id:id},
        function () { },
        function(msg){
            TASKVARS=msg
            $("#task-id").val(TASKVARS.task_id)
            $("#repo-id").val(TASKVARS.repo_id)
			$("#task-name").val(TASKVARS.task_name)
			$('#task-create').removeClass("disabled")
			GetRepoVars(TASKVARS.repo_id,true)
        },
        ReqErr
    );
}

function CleanTask(){
	$("#task-id").val("")
	$("#repo-id").val("")
	$("#task-name").val("")
	$(".inv").empty()
	$("#playbook-parse").empty()
}

function TaskList(){
    var success = function (msg) {
        $("#task-list").empty()
        $.each(msg, function (i, val) {
			var status=`<a href="javascript:RunTask('`+val.task_id+`','`+val.repo_id+`','`+val.task_name+`');" class="btn btn-xs btn-success" id="start" title="启动">
				<i class="icon icon-play"></i>启动
				</a>`
			if (val.task_status=="running"||val.task_status=="waiting"){
				var status=`<a href="javascript:StopTask('`+val.task_id+`');" class="btn btn-xs btn-warning" id="stop" title="停止">
				<i class="icon icon-stop"></i>停止
				</a>`
			}
            var html = `<tr>
            <td>`+ val.task_id + `</td>
			<td>`+ val.task_name + `</td>
			<td>`+ val.task_status + `</td>
            <td>`+ val.created + `</td>
			<td>
				`+status+`
				<a href="javascript:ViewTask('`+val.task_name+`','`+val.task_id+`');" class="btn btn-xs btn-info">查看</a>
				<a href="javascript:CreateTimer('`+val.task_id+`');" class="btn btn-xs btn-primary">创建定时任务</a>
				<a href="javascript:GetTask('`+val.task_id+`');" class="btn btn-xs btn-primary">编辑</a>
                <a href="javascript:DelTask('`+val.task_id+`');" class="btn btn-xs btn-danger">删除</a>
            </td>
        </tr>`
            $('#task-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/task",
        {project_id:project_id},
        function () { },
        success,
        ReqErr
    );
}

function DelTask(id){
	AjaxReq(
        "get",
        "../ansible/common/task/del",
        {task_id:id},
        function () { },
        function(){
			TaskList()
		},
        ReqErr
    );
}

function RunTask(tid,rid,name){
	var success = function (msg) {
        $("#task-tag").empty()
        $.each(msg.tag, function (i, val) {
            var html = `<option value="`+val.tag_value+`">`+ val.tag_name + `</option>`
			$('#task-tag').append(html)
		})
		$("#run-name").val(name)
		$("#run-id").val(tid)
		$("#run-modal").modal("show")
    };
	AjaxReq(
        "get",
        "../ansible/common/vars",
        {repo_id:rid},
        function () { },
        success,
        ReqErr
    );
}

function StartTask(){
	$('#start-task').addClass("disabled")
	AjaxReq(
        "get",
        "../ansible/common/task/start",
        {task_id:$("#run-id").val(),task_tag:$("#task-tag").val()},
        function () { },
        function(){
			setTimeout(function(){
				$('#start-task').removeClass("disabled")
				TaskList()
			},2000)
			ViewTask($("#run-name").val(),$("#run-id").val())
			$("#run-modal").modal("hide")
		},
        function(){
			ReqErr()
			$('#start-task').removeClass("disabled")
		}
    );
}

function StopTask(id){
	AjaxReq(
        "get",
        "../ansible/common/task/stop",
        {task_id:id},
        function () { },
        function(){
			$('#stop').addClass("disabled")
			setTimeout(function(){
				$('#stop').removeClass("disabled")
				TaskList()
				ReqSuccess()
			},2000)
		},
		function(){
			ReqErr()
			$('#stop').removeClass("disabled")
		}
        
    );
}

// 创建定时任务
function CreateTimer(task_id){
	$('#ttask-id').val(task_id)
	$("#timer-modal").modal("show")
}


function CommitTimer(){
	var data = {
		timer_id: $('#timer-id').val(),
		task_id: $('#ttask-id').val(),
        timer_name: $('#timer-name').val(),
        timer_interval: parseInt($('#timer-interval').val()),
        timer_repeat: parseInt($('#timer-repeat').val()),
        timer_status:new Boolean($('#timer-status').val())
        
    }
	AjaxReq(
        "post",
        "../ansible/common/timer/create",
        JSON.stringify(data),
        function () { },
        function(){
			ReqSuccess()
			$(location).attr("href","../ui/timer.html")
		},
        ReqErr
    );
}

function ViewTask(name,id){
	var taskModalTrigger = new $.zui.ModalTrigger({
		title:"任务："+name,
		size:"fullscreen",
		iframe:"term.html?task_id="+id
		});
	taskModalTrigger.show();
}

function CreateTask(){
    $("#playbook-parse").find(".vars-value").each(function(i){
        TASKVARS.task_vars[i].vars_value.vars=JSON.parse($(this).val())
    })
    TASKVARS.task_group=createInvJSON()
    TASKVARS.task_id=$("#task-id").val()
    TASKVARS.repo_id=$("#repo-id").val()
    TASKVARS.task_name=$("#task-name").val()
    TASKVARS.project_id=project_id
	TASKVARS.task_tag=""
    AjaxReq(
        "post",
        "../ansible/common/task/create",
        JSON.stringify(TASKVARS),
        function () { },
        function(msg){
            $('#task-modal').modal('hide');
            TaskList()
            ReqSuccess()
        },
        ReqErr
    );
}


function createInvView(){
	$(".inv").empty()
	$.each(REPOVARS.group,function(i,v){
		var attrHtml=""
		$.each(v.attr,function(i,v){
			var input
			if (v.type=="bool"){
				input=`<div>
					<select class="form-control attr-value">
						<option value="yes">yes</option>
						<option value="no">no</option>
					</select>
				</div>`
				if (v.default=="no"){
					input=`<div>
						<select class="form-control attr-value">
							<option value="yes">yes</option>
							<option value="no" selected>no</option>
						</select>
					</div>`
				}
				
			}else{
				input=`<div><input class="form-control attr-value" type="text" value="`+v.default+`"/></div>`
			}
			attrHtml+=`<div class="inv-attr" style="margin-top:10px;">
							<label  class=" col-md-2 col-sm-6 attr-name">`+v.key+`</label>
							<div class="col-md-10 col-sm-6">
								`+input+`
							</div>
						</div>`
		})
        var host_attr=""
		$.each(HOSTS,function(i,v){
			host_attr+=`<div class="form-group inv-host">
						<div>
								<label>
									<input type="checkbox" class="host-name" value="`+v.host_name+`"> `+v.host_alias+`
								</label>
						</div>`
						+attrHtml+`
						</div>`
		})
		$(".inv").append(`<span class="inv-group">
				<label class="group-name">【`+v.group_name+`】</label>
				`+host_attr+`
			 </span>`)
	})
	onCheck()
}


function editInvView(){
	$(".inv").find(".inv-group").each(function(i){
			$(this).find(".host-name").each(function(j){
				var hostDom=$(this)
				$.each(TASKVARS.task_group[i].hosts,function(k,v){
					if(hostDom.val()==v.host_name){
						hostDom.attr("checked",true)
						hostDom.parents(".inv-host").find(".inv-attr").each(function(m){
							var keyDom=$(this).children(".attr-name")
							var valueDom=$(this).find(".attr-value")
							$.each(TASKVARS.task_group[i].hosts[k].attr,function(n,v){
								if(keyDom.text()==v.key)
									valueDom.val(v.value)
							})
						})
					}
						

				})
				
			})
		})
	check()
}

function createInvJSON(){
	var obj=new Array()
	$(".inv").find(".inv-group").each(function(i){
		var group=new Object()
		group.group_name=$(this).children(".group-name").text()
		group.hosts=new Array()
	   $(this).find(".inv-host").each(function(j){
		   var host=new Object()
		   var select =$(this)
		   $(this).find(".host-name").each(function(i){
				if(this.checked==false){
					return
				}
				host.attr=new Array()
				host.host_name=$(this).val()
				select.find(".inv-attr").each(function(k){
					var key=new Array()
					var value=new Array()
					$(this).find(".attr-name").each(function(i){
						key.push($(this).text())
					})
					$(this).find(".attr-value").each(function(k){
						value.push($(this).val())
					})
					$.each(key,function(i,v){
						var attr ={}
						attr.key=v
						attr.value=value[i]
						host.attr.push(attr)
					})
				})
				group.hosts.push(host)
		   })
	   })
	   obj.push(group)
	})
	return obj
}

function createVarsView(){
	$("#playbook-parse").empty()
	$.each(TASKVARS.task_vars,function(i,v){
		$("#playbook-parse").append( `
		<div class="form-group">
			<label for="playbook-value" class="col-sm-2">`+v.vars_name+`</label>
			<div class="col-md-9 col-sm-10">
				<textarea class="form-control vars-value" rows="10">`+JSON.stringify(v.vars_value.vars).replace(/\\n/g,"\n")+`</textarea>
			</div>
		</div>`)
	})
	$(".inv-attr").hide()
}

function onCheck(){
	$(".inv").find("input[type=checkbox]").each(function(i) {
		$(this).change(function(){
			if(this.checked==true){
				  $(this).parents(".inv-host").find(".inv-attr").show()
			}else{
				$(this).parents(".inv-host").find(".inv-attr").hide()
			}
		 });
	});  
}

function check(){
	$(".inv").find("input[type=checkbox]").each(function(i) {
			if(this.checked==true){
				  $(this).parents(".inv-host").find(".inv-attr").show()
			}else{
				$(this).parents(".inv-host").find(".inv-attr").hide()
			}
	}); 
}
var parseObj=new Object()
parseObj.vars=new Array()
var group_tpl

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

function createInvView(){
	$.each(group_tpl,function(i,v){
		var attrHtml=""
		$.each(v.attr,function(i,v){
			var input
			if (v.type=="bool"){
				input=`<select class="form-control attr-value">
					<option value="yes">yes</option>
					<option value="no">no</option>
				</select>`
				if (v.default=="no"){
					input=`<select class="form-control attr-value">
						<option value="yes">yes</option>
						<option value="no" selected>no</option>
					</select>`
				}
				
			}else{
				input=`<input class="form-control attr-value" type="text" value="`+v.default+`"/>`
			}
			attrHtml+=`<div class="col-md-4 col-sm-10 inv-attr" >
							<label  class="col-sm-6 attr-name">`+v.key+`</label>
							<div class="col-sm-6">
								`+input+`
							</div>
						</div>`
		})
		var host_attr=""
		$.each(parseObj.hosts,function(i,v){
			host_attr+=`<div class="form-group inv-host">
						<div class="col-sm-2">
								<label>
									<input type="checkbox" class="host-name" value="`+v.host_name+`"> `+v.host_alias+`
								</label>
						</div>`
						+attrHtml+`
						</div>`
		})
		$(".inv").append(`<span class="inv-group">
				<label class="group-name">`+v.group_name+`</label>
				`+host_attr+`
			 </span>`)
	})
	onCheck()
}

function createVarsView(){
	$.each(parseObj.vars,function(i,v){
		$("#playbook-parse").append( `
		<div class="form-group">
			<label for="playbook-value" class="col-sm-2">`+v.vars_name+`</label>
			<div class="col-md-9 col-sm-10">
				<textarea class="form-control vars-value" rows="10">`+v.vars_value.replace(/\\n/g,"\n")+`</textarea>
			</div>
		</div>`)
	})
}

function editInvView(){
	obj=parseObj.group
	$(".inv").find(".inv-group").each(function(i){
			$(this).find(".host-name").each(function(j){
				var hostDom=$(this)
				$.each(obj[i].hosts,function(k,v){
					if(hostDom.val()==v.host_name){
						hostDom.attr("checked",true)
						hostDom.parents(".inv-host").find(".inv-attr").each(function(m){
							var keyDom=$(this).children(".attr-name")
							var valueDom=$(this).find(".attr-value")
							$.each(obj[i].hosts[k].attr,function(n,v){
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
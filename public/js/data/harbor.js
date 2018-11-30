
$(function(){
    GetRepoList()
})



// get host list
function GetRepoList() {
    var success = function (msg) {
        $("#repo-list").empty()
        $.each(msg, function (i, val) {
           var versions =`<select  id="version-`+i+`">`
            $.each(val.charts,function(i,v){
                versions  +=`<option value="`+v.version+`">`+ v.appVersion+"/"+v.version + `</option>`
            })
            versions  +=`<option value=""> latest</option>`
            versions +="</select>"
       
        var html =`<div class="col-md-4 col-sm-6 col-lg-3">
        <span class="card">
            <a href="###" style="margin:10px;">
                <p style="margin: 0 auto; width: 120px; height: 120px">
                    <img width="100%" height="100%"onerror='this.src="img/default.png"' src="`+val.icon+`" title="点击部署">
                </p>
            </a>
            <br/>
            <br/>
            <div class="card-heading">
                    <strong>`+ val.name + `</strong>
                   `+versions+`
            </div>
            <div class="card-content text-muted" style="height:60px;">`+ val.description + `</div>
            <div class="card-actions">
               <a href="javascript:Release('`+val.name+`','`+i+`')" type="button" class="btn"><i class="icon-play-circle"></i> 发布</a>
            </div>
        </span>
    </div>`
            $('#repo-list').append(html)
        })
    };

    AjaxReq(
        "get",
        "../ansible/common/harbor/charts",
        {},
        function () { },
        success,
        ReqErr
    );
}


function Release(chart,i){
    var version =$("#version-"+i).val()
    console.log(version)
    window.location.href="helm-install.html?chart="+chart+"&version="+version;
}
$('#login').click(function(){
    var account=$("#username").val()
    var password=$("#password").val()
    $.ajax({
        type: "post",
        url: "../ansible/login",
        data: {
            account:account,
            password:password
        },
        beforeSend: function() {},
        success: function(msg) {
            $.cookie("Auth",msg.token,{path:'/'})
            $.cookie("Auth",msg.token)
            $(location).attr("href","../ui/")
        },
        error: function(XMLHttpRequest, textStatus, errorThrown) {
            new $.zui.Messager('提示消息：登录失败', {
                type: 'danger' 
            }).show();
        }
    });
})


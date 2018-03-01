package routers

import (
	"net/http"

	"github.com/gzsunrun/ansible-manager/controllers"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/auth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
	}))

	authApi := beego.NewNamespace("/ansible",
		beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
	)

	commonApi:=beego.NewNamespace("/ansible/common",
		beego.NSRouter("/user/info", &controllers.UserController{}, "get:Info"),
		beego.NSRouter("/user/create", &controllers.UserController{}, "post:Create"),
		beego.NSRouter("/user", &controllers.UserController{}, "get:List"),
		beego.NSRouter("/user/del", &controllers.UserController{}, "get:Del"),
		beego.NSRouter("/repo", &controllers.RepoController{}, "get:List"),
		beego.NSRouter("/repo/create", &controllers.RepoController{}, "post:Create"),
		beego.NSRouter("/repo/delete", &controllers.RepoController{}, "get:Delete"),
		beego.NSRouter("/vars", &controllers.RepoController{}, "get:Vars"),
		beego.NSRouter("/hosts", &controllers.HostController{}, "get:List"),
		beego.NSRouter("/hosts/create", &controllers.HostController{}, "post:Create"),
		beego.NSRouter("/hosts/get", &controllers.HostController{}, "get:Get"),
		beego.NSRouter("/hosts/del", &controllers.HostController{}, "get:Del"),
		beego.NSRouter("/project", &controllers.ProjectController{}, "get:List"),
		beego.NSRouter("/project/get", &controllers.ProjectController{}, "get:GetProject"),
		beego.NSRouter("/project/create", &controllers.ProjectController{}, "post:Create"),
		beego.NSRouter("/project/caa", &controllers.ProjectController{}, "post:CreateAndAddHosts"),
		beego.NSRouter("/project/del", &controllers.ProjectController{}, "get:Del"),
		beego.NSRouter("/project/addhost", &controllers.ProjectController{}, "post:AddHost"),
		beego.NSRouter("/project/delhost", &controllers.ProjectController{}, "post:DelHost"),
		beego.NSRouter("/project/hosts", &controllers.ProjectController{}, "get:HostList"),
		beego.NSRouter("/task", &controllers.TaskController{}, "get:List"),
		beego.NSRouter("/task/create", &controllers.TaskController{}, "post:Create"),
		beego.NSRouter("/task/start", &controllers.TaskController{}, "get:Start"),
		beego.NSRouter("/task/stop", &controllers.TaskController{}, "get:Stop"),
		beego.NSRouter("/task/get", &controllers.TaskController{}, "get:Get"),
		beego.NSRouter("/task/del", &controllers.TaskController{}, "get:Del"),
		beego.NSRouter("/task/notes", &controllers.TaskController{}, "get:GetNotes"),
		beego.NSRouter("/task/count", &controllers.TaskController{}, "get:GetTaskCount"),
	)

	iaasApi := beego.NewNamespace("/api/ansible",
		beego.NSRouter("/repo", &controllers.RepoController{}, "get:List"),
		beego.NSRouter("/repo/create", &controllers.RepoController{}, "post:Create"),
		beego.NSRouter("/repo/delete", &controllers.RepoController{}, "get:Delete"),
		beego.NSRouter("/vars", &controllers.IaaSController{}, "get:RepoVars"),
		beego.NSRouter("/host", &controllers.IaaSController{}, "get:ProjectHosts"),
		beego.NSRouter("/host/update", &controllers.IaaSController{}, "post:CreateHost"),
		beego.NSRouter("/task/create", &controllers.IaaSController{}, "post:CreateTask"),
		beego.NSRouter("/task/history", &controllers.IaaSController{}, "get:GetTask"),
		beego.NSRouter("/task/stop", &controllers.IaaSController{}, "get:Stop"),
		beego.NSRouter("/task/notes", &controllers.TaskController{}, "get:GetNotes"),
	)

	beego.AddNamespace(authApi,iaasApi,commonApi)
	beego.InsertFilter("/ansible/common/*", beego.BeforeRouter, auth.JwtAuthFilter)
	beego.InsertFilter("/api/ansible/*", beego.BeforeRouter, auth.IaaSAuthFilter)

	beego.Handler("/api/ansible/ws", socketHandler(sockets.Handler))
}

type socketHandler func(http.ResponseWriter, *http.Request)

func (fn socketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fn(w, r)
}

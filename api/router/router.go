package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gzsunrun/ansible-manager/api/project"
	"github.com/gzsunrun/ansible-manager/api/sockets"
)

func NewRouter(root *mux.Router) {
	dir := "/var/lib/ansible-manager"
	root.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(dir+"/public/"))))
	root.HandleFunc("/ws", sockets.Handler)
	root.HandleFunc("/api/login", project.Login)
	projectRoute := root.PathPrefix("/api/project").Subrouter()
	projectRoute.Handle("", apiHandler(project.GetProject))
	projectRoute.Handle("/get", apiHandler(project.GetProjectByID))
	projectRoute.Handle("/create", apiHandler(project.CreateProject))
	projectRoute.Handle("/update", apiHandler(project.UpdateProject))
	projectRoute.Handle("/delete", apiHandler(project.DeleteProject))
	userRoute := root.PathPrefix("/api/user").Subrouter()
	userRoute.Handle("", apiHandler(project.GetUser))
	userRoute.Handle("/current", apiHandler(project.GetCurrentUser))
	userRoute.Handle("/get", apiHandler(project.GetUserByID))
	userRoute.Handle("/create", apiHandler(project.AddUser))
	userRoute.Handle("/update", apiHandler(project.UpdateUser))
	repoRoute := root.PathPrefix("/api/repo").Subrouter()
	repoRoute.Handle("", apiHandler(project.GetRepository))
	repoRoute.Handle("/get", apiHandler(project.GetRepositoryID))
	repoRoute.Handle("/create", apiHandler(project.CreateRepository))
	repoRoute.Handle("/delete", apiHandler(project.DeleteRepository))
	repoRoute.Handle("/git", apiHandler(project.CloneGitRepo))
	varsRoute := root.PathPrefix("/api/repo/vars").Subrouter()
	varsRoute.Handle("", apiHandler(project.GetVars))
	varsRoute.Handle("/get", apiHandler(project.GetVarsByID))
	varsRoute.Handle("/create", apiHandler(project.CreateVars))
	varsRoute.Handle("/delete", apiHandler(project.DeleteVars))
	varsRoute.Handle("/update", apiHandler(project.UpdateVars))
	varsRoute.Handle("/tag", apiHandler(project.GetTagByTplID))
	templateRoute := root.PathPrefix("/api/project/template").Subrouter()
	templateRoute.Handle("", apiHandler(project.GetTemplate))
	templateRoute.Handle("/get", apiHandler(project.GetTemplateByID))
	templateRoute.Handle("/create", apiHandler(project.CreateTemplate))
	templateRoute.Handle("/update", apiHandler(project.UpdateTemplate))
	templateRoute.Handle("/delete", apiHandler(project.DeleteTemplate))
	taskRoute := root.PathPrefix("/api/project/task").Subrouter()
	taskRoute.Handle("", apiHandler(project.GetTask))
	taskRoute.Handle("/get", apiHandler(project.GetTaskByID))
	taskRoute.Handle("/create", apiHandler(project.CreateTask))
	taskRoute.Handle("/delete", apiHandler(project.DeleteTask))
	taskRoute.Handle("/output", apiHandler(project.GetHistory))
	taskRoute.Handle("/stop", apiHandler(project.StopTask))
	hostRoute := root.PathPrefix("/api/project/host").Subrouter()
	hostRoute.Handle("", apiHandler(project.GetHost))
	hostRoute.Handle("/get", apiHandler(project.GetHostByID))
	hostRoute.Handle("/create", apiHandler(project.CreateHost))
	hostRoute.Handle("/update", apiHandler(project.UpdateHost))
	hostRoute.Handle("/delete", apiHandler(project.DeleteHost))
}

type apiHandler func(http.ResponseWriter, *http.Request)

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := project.Auth(w, r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fn(w, r)
}

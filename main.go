package main

import (
    "net/http"
    "fmt"
    "github.com/weekndCN/jenkinsAPI/pkg/jenkins"
)

// GetJobs return REST API record to client
func GetJobs(w http.ResponseWriter, r *http.Request)  {
    w.Header().Set("Access-Control-Allow-Headers", "*")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    auth := &jenkins.Auth{
		Username: "weeknd",
		APIToken: "1166c439b661e415ade72bc6d3fcff4211",
	}
	jenkins := jenkins.NewJenkins(auth, "http://t01.corp.wukongbox.cn:9090")
    jobs, err := jenkins.GetJobs(2)
    if err == nil {
        fmt.Fprint(w, jobs)
    }
}

// GetBuild return job buid log
func GetBuild(w http.ResponseWriter, r *http.Request)  {
    buildurl := r.URL.Query().Get("buildurl")
    fmt.Println(buildurl)
    w.Header().Set("Access-Control-Allow-Headers", "*")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    auth := &jenkins.Auth{
		Username: "weeknd",
		APIToken: "1166c439b661e415ade72bc6d3fcff4211",
	}
    // jenkins := jenkins.NewJenkins(auth, "http://t01.corp.wukongbox.cn:9090")
    //const url = "http://t01.corp.wukongbox.cn:9090/job/scale-backend/141"
    jenkins := jenkins.NewJenkins(auth, "")
    buildlog, err := jenkins.GetBuild(buildurl)
    if err == nil {
        fmt.Fprint(w, buildlog)
    }
}

// HTTP Handler
func main() {
    http.HandleFunc("/getjobs", GetJobs)
    http.HandleFunc("/getbuild", GetBuild)
    http.ListenAndServe(":8080", nil)
}
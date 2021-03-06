package jenkins

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

// Auth is jenkins api token certified
type Auth struct {
	Username string 
	APIToken string
}

// Jenkins full URI
type Jenkins struct {
	auth 	*Auth
	baseURL string
	client 	*http.Client
}

// NewJenkins Returns Jenkins address
func NewJenkins(auth *Auth, baseURL string) *Jenkins {
	return &Jenkins {
		auth: 		auth,
		baseURL: 	baseURL,
		client: 	http.DefaultClient,
	}
}

// get func new request to jenkins
func (jenkins *Jenkins) get(path, params string, body *interface{}, depth int) (err error) {
	requestURL := jenkins.buildURL(path, params, depth)
	req, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return
	}

	resp, err := jenkins.sendRequest(req)
	if err != nil {
		return
	}

	return jenkins.parseResponse(resp, body)
}

// buildURL to get build details
func (jenkins *Jenkins) buildURL(path, params string, depth int) (requestURL string) {
	/*
	requestURL = jenkins.baseURL + "/api/json?" + path
	if params != "" {
		requestURL = jenkins.baseURL + params + "/api/json?" + path
	}
	*/
	switch params {
	// eg: http://xxxx/job/job_name/111/consoleText
	case "log":
		requestURL = path + "consoleText"
	case "jobs":
		requestURL =  jenkins.baseURL + "/api/json?"
	default:
		requestURL =  jenkins.baseURL + "/api/json?"
	}

	if depth > 0 {
		requestURL = requestURL + fmt.Sprintf("&depth=%d",depth)
	} 

	fmt.Println(requestURL)

	return
}


// sendRequest do request
func (jenkins *Jenkins) sendRequest(req *http.Request) (*http.Response, error) {
	// add username and APIToken to header	
	if jenkins.auth != nil {
		req.SetBasicAuth(jenkins.auth.Username, jenkins.auth.APIToken)
	}
	return jenkins.client.Do(req)
}

// parseResponse to parse response body
func (jenkins *Jenkins) parseResponse(resp *http.Response, body *interface {}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	
	*body = string(data)

	return
}

// GetJobs is jenkins methods
func (jenkins *Jenkins) GetJobs(depth int) (interface{}, error) {
	var jobs interface{}
	err := jenkins.get("", "", &jobs, depth)
	return jobs, err
}

// GetBuild to get build log of a specified job name
func (jenkins *Jenkins) GetBuild(url string) (interface{}, error)  {
	var buildlog interface{}
	err := jenkins.get(url, "log", &buildlog, 0)
	return buildlog, err
}
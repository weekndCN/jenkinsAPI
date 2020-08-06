package jenkins

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"io/ioutil"
)

// define const value
const (
	jobsPath 	= "tree=jobs[displayName,description,displayNameOrNull,fullDisplayName,fullName,name,url]"
	viewPath 	= "tree=views[description,name,url,jobs]"
	buildPath 	= "tree=allBuilds[description,displayName,duration,estimatedDuration,result,timestamp,id,number,url,builds]&depth=2"
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

// GetJob to get details by a specified name of job
func (jenkins *Jenkins) GetJob(name string) (job interface{}, err error) {
	err = jenkins.get(fmt.Sprintf("/job/%s", name), nil, &job)
	return
}

// get func new request to jenkins
func (jenkins *Jenkins) get(path string, params url.Values, body interface{}) (err error) {
	requestURL := jenkins.buildURL(path, params)
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
func (jenkins *Jenkins) buildURL(path string, params url.Values) (requestURL string) {
	requestURL = jenkins.baseURL + path + "/api/json"
	if params != nil {
		queryString := params.Encode()
		if queryString != "" {
			requestURL = requestURL + "?" + queryString
		}
	}

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
func (jenkins *Jenkins) parseResponse(resp *http.Response, body interface{}) (err error) {
	defer resp.Body.Close()

	if body == nil {
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, body)
}

// GetJobs is jenkins methods
func (jenkins *Jenkins) GetJobs() (interface{}, error) {
	var payload interface{}
	err := jenkins.get("", nil, &payload)
	return payload, err
}
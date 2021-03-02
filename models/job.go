package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AjouJob 아주대학교 채용정보
var AjouJob = os.Getenv("AJOU_JOB")

// Job ...
type Job struct {
	Code string `json:"msgCode"`
	Data []struct {
		Title string `json:"noti"`
		Date  string `json:"date"`
		URL   string `json:"url"`
	} `json:"p09List"`
}

// GetJobsAvailable 교내 채용 정보를 10개 불러옴
func GetJobsAvailable() (Job, error) {
	var job Job

	// POST
	jsonValue, _ := json.Marshal(map[string]int{"tabIndex": 0})

	// Disable SSL authentication and post jsonValue
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post((AjouJob), "application/json;charset=UTF-8", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
		return job, err
	}
	defer resp.Body.Close()

	// Response 체크.
	respBody, err := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &job)
	if err != nil {
		fmt.Println(err)
		return job, err
	}

	job.Data = job.Data[:10]

	return job, nil
}

package pipelinerun

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetPipelineRun(host string, job string, runId int, username string, password string) PipelineRun {
	// https://jenkins.gke.kearos.net/blue/rest/organizations/jenkins/pipelines/prom-test-2/runs/7/nodes/

	url := fmt.Sprintf("%v/blue/rest/organizations/jenkins/pipelines/%v/runs/%v/nodes/", host, job, runId)

	fmt.Printf("> Retrieving Pipeline Run from: %v\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	var pipelineRun PipelineRun
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	json.Unmarshal([]byte(body), &pipelineRun)
	fmt.Printf(" > Found %v stages in the pipeline run\n", len(pipelineRun))
	for _, stage := range pipelineRun {
		fmt.Printf(" - Stage: %v (%v)\n", stage.DisplayName, stage.ID)
		fmt.Printf("   [ %v, %v, %v, %v]\n", stage.State, stage.Result, stage.Type, stage.DurationInMillis)
	}
	return pipelineRun
}

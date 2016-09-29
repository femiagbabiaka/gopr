package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// PullRequest represents a the structure of our JSON body for Github.
type PullRequest struct {
	Title       string `json:"title"`
	Target      string `json:"head"`
	Base        string `json:"base"`
	Description string `json:"body"`
	requestURL  string
}

// TargetRepositoryInfo represents the repo section of our API endpoint for Github.
type TargetRepositoryInfo struct {
	Owner      string
	Repository string
}

func main() {

	var title, target, base, description, targetRepository string

	flag.StringVar(&title, "title", "", "The title of your pull request.")
	flag.StringVar(&target, "target", "", "The target branch of your pull request.")
	flag.StringVar(&base, "base", currentBranch(), "The base branch for your pull request.")
	flag.StringVar(&description, "description", "", "The description of your pull request.")
	flag.StringVar(&targetRepository, "targetrepository", "", "The target repository for your pull request.")

	flag.Parse()

	if title == "" {
		log.Fatal("Must supply a pull request title.")
	}

	if target == "" {
		log.Fatal("Must supply a pull request target.")
	}

	if base == "" {
		log.Fatal("Must supply a base branch for your pull request.")
	}

	if description == "" {
		log.Fatal("Must supply a description for your pull request.")
	}

	if targetRepository == "" {
		log.Fatal("Must supply a target repository for your pull request.")
	}

	requestRepositoryInfo := generateTargetRepositoryInfo(targetRepository)
	requestURL := requestRepositoryInfo.createRequestURL()

	pr := PullRequest{Title: title, Target: target, Base: base, Description: description, requestURL: requestURL}

	res := pr.pullRequest()

	if res {
		os.Exit(0)
	}

	os.Exit(1)
}

func generateTargetRepositoryInfo(tarRep string) TargetRepositoryInfo {
	compositeStrings := strings.Split(tarRep, "/")

	return TargetRepositoryInfo{Owner: compositeStrings[0], Repository: compositeStrings[1]}
}

// TODO: SUPPORT ENVIRONMENT VARIABLES FOR CONFIG STUFF
func (r TargetRepositoryInfo) createRequestURL() string {
	return fmt.Sprintf("/repos/%s/%s/pulls", r.Owner, r.Repository)
}

func (p PullRequest) pullRequest() bool {
	client := &http.Client{Timeout: 0}

	requestBody, err := json.Marshal(p)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com%s", p.requestURL), bytes.NewBuffer(requestBody))
	//TODO: Set Headers here, read auth token and OAuth from Environment variables.

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if body != nil {
		return true
	}

	return false
}

func currentBranch() string {
	out, err := exec.Command("git branch | grep * | cut -d ' ' -f2").Output()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(out[:])
}

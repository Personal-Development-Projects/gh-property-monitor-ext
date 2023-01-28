package main

import (
	"encoding/json"
	"fmt"
	"github.com/cli/go-gh"
	"strconv"
)

func main() {

	// Using github rest api to get username that is logged in at enterprise server
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Starting gh-property-monitor extension")
	fmt.Printf("Running as %s\n", response.Login)

	fmt.Println("Collecting PRs associated with branch")
	//TODO Would like to have jq return more formatted return in order to reduce parsing
	updatedPRBuffer, _, err := gh.Exec("search", "prs", "--merged-at", "", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "--json", "number,repository,author")
	if err != nil {
		fmt.Println(err)
		return
	}
	var prBaseResults PullRequestsResults
	err = json.Unmarshal(updatedPRBuffer.Bytes(), &prBaseResults)
	if err != nil {
		fmt.Println(err)
	}
	// Function to populate these four arrays of struct
	// 		FileNameAdditions []string
	//		FileNameDeletions []string
	//		textAdditions     []string
	//		textDeletions     []string
	populateDetailedResults(prBaseResults)

	for i := 0; i < len(prBaseResults); i++ {
	}
	//TODO Incorporate unmarshalled object into PR details function
	//fmt.Println("Starting to collect the file changes from PR's")
}

func populateDetailedResults(prList PullRequestsResults) {
	// Iterate through the list of PRs
	for prIndex := 0; prIndex < len(prList); prIndex++ {
		prNumString := strconv.Itoa(prList[prIndex].Number)
		prDetails, _, err := gh.Exec("pr", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "diff", prNumString)
		if err != nil {
			fmt.Println(err)
		}
		// parse those into struct that can be merged into prList
		parseDetailedArrays(prList[prIndex], prDetails.String())
	}
}

func parseDetailedArrays(pr struct{ PullRequest }, prDetailsRaw string) {
	fmt.Println(pr)
}

type PullRequestsResults []struct {
	PullRequest
}

type PullRequest struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	Number     int `json:"number"`
	Repository struct {
		Name string `json:"name"`
		//NameWithOwner string `json:"nameWithOwner"`
	} `json:"repository"`
	DetailedResults struct {
		FileNameAdditions []string
		FileNameDeletions []string
		textAdditions     []string
		textDeletions     []string
	}
}

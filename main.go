package main

import (
	"fmt"
	"github.com/cli/go-gh"
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
	fmt.Println("Starting to collect the file changes from PR's")
	updatedPRString := updatedPRBuffer.String()
	fmt.Println(updatedPRString)

}

//func getPRDetails(prNumber string) {
//	prDetails, _, err := gh.Exec("pr", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "diff", prNumber)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Print(prDetails.String())
//	//fmt.Println(_prDetails)
//}

type name struct {
}

type PullRequest struct {
	Number     string
	Repository string
	Author     string
	//Approver string
	FileNameAdditions []string
	FileNameDeletions []string
	textAdditions     []string
	textDeletions     []string
}

/*
gh-property-monitor is a gh cli extension that can be utilized to monitor PR's of an associated Repo

It uses the GH CLI tool to access online GitHub repositories and collect Pull Request information for logging in
order to reduce manual intervention of property file validation.

Usage:

1) Navigate to the directory in which this program resides on your local machine
2) Ensure you are logged into gh cli by using ```gh auth status```
3) Once authorized, install extension with ```gh extension install gh-property-monitor``` or ```gh extension install .```
4) To run the extension use ``` gh property-monitor```

The output from running this program will be within an Excel file named Property-Monitor and the results for extension run
can be found under today's local date
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/cli/go-gh"
	"strconv"
	"strings"
)

// This is a test comment
func main() {

	// Using GitHub rest api to get username that is logged in at remote server
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

	//TODO Possibly add user option to log in or verify that this username is correct
	fmt.Println("Collecting PRs associated with branch")
	updatedPRBuffer, _, err := gh.Exec("search", "prs", "--merged-at", "", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "--json", "number,author")
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
	finalResult := populateDetailedResults(prBaseResults)

	// ExcelUtility.go Output the results to Excel file
	writeResultsToFile(finalResult)
}

func populateDetailedResults(prList PullRequestsResults) PullRequestsResults {
	// Iterate through the list of PRs
	for prIndex := 0; prIndex < len(prList); prIndex++ {
		prNumString := strconv.Itoa(prList[prIndex].Number)
		prDetails, _, err := gh.Exec("pr", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "diff", prNumString)
		if err != nil {
			fmt.Println(err)
		}
		// parse those into struct that can be merged into prList
		scanner := bufio.NewScanner(&prDetails)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			resultString := scanner.Text() + ", \n"
			if strings.HasPrefix(scanner.Text(), "+") {
				if strings.Count(scanner.Text(), "+") >= 3 {
					prList[prIndex].FileAdditions += resultString
				} else {
					prList[prIndex].TxtAdditions += resultString
				}
			} else if strings.HasPrefix(scanner.Text(), "-") {
				if strings.Count(scanner.Text(), "-") >= 3 {
					prList[prIndex].FileDeletions += resultString
				} else {
					prList[prIndex].TxtDeletions += resultString
				}
			}

		}
	}
	return prList
}

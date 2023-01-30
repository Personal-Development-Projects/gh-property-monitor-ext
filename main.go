package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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

	//TODO Possibly add user option to login or verify that this username is correct
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
	var finalResult = populateDetailedResults(prBaseResults)
	writeResultsToFile(finalResult)

}

// TODO This could be improved by increasing the proficiency of the parser
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
			if strings.HasPrefix(scanner.Text(), "+") {
				if strings.Count(scanner.Text(), "+") >= 3 {
					prList[prIndex].DetailedResults.FileAdditions = append(prList[prIndex].DetailedResults.FileAdditions, scanner.Text())
				} else {
					prList[prIndex].DetailedResults.TxtAdditions = append(prList[prIndex].DetailedResults.TxtAdditions, scanner.Text())
				}
			}

			if strings.HasPrefix(scanner.Text(), "-") {
				if strings.Count(scanner.Text(), "-") >= 3 {
					prList[prIndex].DetailedResults.FileDeletions = append(prList[prIndex].DetailedResults.FileDeletions, scanner.Text())
				} else {
					prList[prIndex].DetailedResults.TxtDeletions = append(prList[prIndex].DetailedResults.TxtDeletions, scanner.Text())
				}
			}

		}
	}
	return prList
}

// Need to iterate through results and write to file in correct format
func writeResultsToFile(resultList PullRequestsResults) {
	fileName := time.Now().Format("2006-01-02")
	fmt.Println(fileName + ".csv")
	file, err := os.Create(fileName + ".csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	createCSVDataArray(resultList)

}

func createCSVDataArray(requestsResults PullRequestsResults) {
	resultsCSV := [][]string{{}}
	for index := 0; index < len(requestsResults); index++ {

	}
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
	DetailedResults
}

type DetailedResults struct {
	FileAdditions []string `json:"_"`
	TxtAdditions  []string `json:"_"`
	FileDeletions []string `json:"_"`
	TxtDeletions  []string `json:"_"`
}

func (resultDetails DetailedResults) String() string {
	for i := 0; i < 4; i++ {

	}
}

func (pr PullRequest) String() []string {
	resultString := []string{string(pr.Number), pr.Author.Login, pr.DetailedResults.String()}
	return resultString
}

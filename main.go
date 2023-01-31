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
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"time"
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

	// Output the results to Excel file
	writeResultsToFile(finalResult)
}

func writeResultsToFile(results PullRequestsResults) {
	// Create new file
	//TODO make this load existing file or create a new one if needed
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet to work with
	sheetName := time.Now().Format("2006-01-02")
	currentSheet, err := file.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Setup Title
	file.MergeCell(sheetName, "A1", "F1")
	sheetTitle := "MGMT Config Query - " + sheetName
	file.SetCellValue(sheetName, "A1", sheetTitle)
	// Set headers
	file.SetCellValue(sheetName, "A2", "PR Number")
	file.SetCellValue(sheetName, "B2", "Author")
	file.SetCellValue(sheetName, "C2", "Files Changed (Additions)")
	file.SetCellValue(sheetName, "D2", "Line Additions")
	file.SetCellValue(sheetName, "E2", "Files Changed (Deletions)")
	file.SetCellValue(sheetName, "F2", "Line Deletions")

	populateDataInExcel(file, sheetName, results)

	file.SetActiveSheet(currentSheet)

	if err := file.SaveAs("test1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func populateDataInExcel(file *excelize.File, sheetName string, results PullRequestsResults) {
	for i := 0; i < len(results); i++ {
		//fmt.Println("PR Number: ", csvReadyResults[i].Number, " Author: ", csvReadyResults[i].Author.Login, " Files Affected (additions): ", csvReadyResults[i].FileAdditions, " Text Additions: ", csvReadyResults[i].TxtAdditions, " Files Affected (deletions): ", csvReadyResults[i].FileDeletions, " Text Deletions: ", csvReadyResults[i].TxtDeletions)
		cellNumber := i + 3
		putPRNumberInTable(file, sheetName, cellNumber, results[i].Number)
		putAuthorInTable(file, sheetName, cellNumber, results[i].Author.Login)
		putFileAdditionsInTable(file, sheetName, cellNumber, results[i].FileAdditions)
		putTxtAdditionsInTable(file, sheetName, cellNumber, results[i].TxtAdditions)
		putFileDeletionsInTable(file, sheetName, cellNumber, results[i].FileDeletions)
		putTxtDeletionsInTable(file, sheetName, cellNumber, results[i].TxtDeletions)
	}
}

func putPRNumberInTable(file *excelize.File, sheetName string, cellNumber int, prNum int) {
	var cellCord = "A" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, prNum); err != nil {
		fmt.Println(err)
	}
}

func putAuthorInTable(file *excelize.File, sheetName string, cellNumber int, login string) {
	var cellCord = "B" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, login); err != nil {
		fmt.Println(err)
	}
}

func putFileAdditionsInTable(file *excelize.File, sheetName string, cellNumber int, fileAdditions string) {
	var cellCord = "C" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, fileAdditions); err != nil {
		fmt.Println(err)
	}
}
func putTxtAdditionsInTable(file *excelize.File, sheetName string, cellNumber int, txtAdditions string) {
	var cellCord = "D" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, txtAdditions); err != nil {
		fmt.Println(err)
	}
}

func putFileDeletionsInTable(file *excelize.File, sheetName string, cellNumber int, fileDeletions string) {
	var cellCord = "E" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, fileDeletions); err != nil {
		fmt.Println(err)
	}
}
func putTxtDeletionsInTable(file *excelize.File, sheetName string, cellNumber int, txtDeletions string) {
	var cellCord = "F" + strconv.Itoa(cellNumber)
	if err := file.SetCellValue(sheetName, cellCord, txtDeletions); err != nil {
		fmt.Println(err)
	}
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
	FileAdditions string `json:"file-additions"`
	TxtAdditions  string `json:"txt-additions"`
	FileDeletions string `json:"file-deletions"`
	TxtDeletions  string `json:"txt-deletions"`
}

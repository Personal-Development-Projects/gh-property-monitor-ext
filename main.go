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
	finalResult := populateDetailedResults(prBaseResults)
	//marshalResultsToCSVExcel(finalResult)

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
	index, err := file.NewSheet(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set headers
	file.SetCellValue(sheetName, "A1", "PR Number")
	file.SetCellValue(sheetName, "B1", "Author")
	file.SetCellValue(sheetName, "C1", "Files Changed (Additions)")
	file.SetCellValue(sheetName, "D1", "Line Additions")
	file.SetCellValue(sheetName, "E1", "Files Changed (Deletions)")
	file.SetCellValue(sheetName, "F1", "Line Deletions")

	// testing output to excel basic
	for i := 0; i < len(results); i++ {
		cellNumber := i + 2
		// Add PRs
		var cellCordA = "A" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordA, results[i].Number)
		// Add Author
		var cellCordB = "B" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordB, results[i].Author.Login)
		// Add file names that are affected in additions
		var cellCordC = "C" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordC, results[i].FileAdditions)
		// Add text additions
		var cellCordD = "D" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordD, results[i].TxtAdditions)
		// Add file names that are affected in deletions
		var cellCordE = "E" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordE, results[i].FileDeletions)
		// Add text deletions
		var cellCordF = "F" + strconv.Itoa(cellNumber)
		file.SetCellValue(sheetName, cellCordF, results[i].TxtDeletions)
	}

	//populateDataInExcel(file, sheetName, results)

	file.SetActiveSheet(index)

	if err := file.SaveAs("test1.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func populateDataInExcel(file *excelize.File, sheetName string, results PullRequestsResults) {
	for i := 0; i < len(results); i++ {
		//fmt.Println("PR Number: ", csvReadyResults[i].Number, " Author: ", csvReadyResults[i].Author.Login, " Files Affected (additions): ", csvReadyResults[i].FileAdditions, " Text Additions: ", csvReadyResults[i].TxtAdditions, " Files Affected (deletions): ", csvReadyResults[i].FileDeletions, " Text Deletions: ", csvReadyResults[i].TxtDeletions)
		cellNumber := i + 2
		putPRNumberInTable(file, sheetName, cellNumber, results[i].Number)
		putAuthorInTable(file, sheetName, cellNumber, results[i].Author.Login)
		putFileAdditionsInTable(file, sheetName, cellNumber, results[i].FileAdditions)
		putTxtAdditionsInTable(file, sheetName, cellNumber, results[i].TxtAdditions)
		putFileDeletionsInTable(file, sheetName, cellNumber, results[i].FileDeletions)
		putTxtDeletionsInTable(file, sheetName, cellNumber, results[i].TxtDeletions)
	}
}

func putPRNumberInTable(file *excelize.File, sheetName string, cellNumber int, prNum int) {
	var cellCord = "A" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, prNum)
}

func putAuthorInTable(file *excelize.File, sheetName string, cellNumber int, login string) {
	var cellCord = "B" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, login)
}

func putFileAdditionsInTable(file *excelize.File, sheetName string, cellNumber int, fileAdditions string) {
	var cellCord = "C" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, fileAdditions)
}
func putTxtAdditionsInTable(file *excelize.File, sheetName string, cellNumber int, txtAdditions string) {
	var cellCord = "D" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, txtAdditions)
}

func putFileDeletionsInTable(file *excelize.File, sheetName string, cellNumber int, fileDeletions string) {
	var cellCord = "E" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, fileDeletions)
}
func putTxtDeletionsInTable(file *excelize.File, sheetName string, cellNumber int, txtDeletions string) {
	var cellCord = "F" + string(cellNumber)
	file.SetCellValue(sheetName, cellCord, txtDeletions)
}

func marshalResultsToCSVExcel(results PullRequestsResults) {
	rawJson, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}
	var csvReadyResults PullRequestsResults
	err = json.Unmarshal(rawJson, &csvReadyResults)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(csvReadyResults)
	for i := 0; i < len(csvReadyResults); i++ {
		fmt.Println("PR Number: ", csvReadyResults[i].Number, " Author: ", csvReadyResults[i].Author.Login, " Files Affected (additions): ", csvReadyResults[i].FileAdditions, " Text Additions: ", csvReadyResults[i].TxtAdditions, " Files Affected (deletions): ", csvReadyResults[i].FileDeletions, " Text Deletions: ", csvReadyResults[i].TxtDeletions)

	}
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

type CSVReadyResults []struct {
	Number int `json:"number"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	FileAdditions string `json:"file-additions"`
	TxtAdditions  string `json:"txt-additions"`
	FileDeletions string `json:"file-deletions"`
	TxtDeletions  string `json:"txt-deletions"`
}

//func (resultDetails DetailedResults) String() string {
//	detailsCSVString := ""
//	//for i := 0; i < 4; i++ {
//	//	detailsCSVString += resultDetails[i]
//	//}
//	return detailsCSVString
//}
//
//func (pr PullRequest) String() []string {
//	resultString := []string{string(pr.Number), pr.Author.Login, pr.DetailedResults.String()}
//	return resultString
//}

package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

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

package main

import (
	"bufio"
	"fmt"
	"github.com/cli/go-gh"
	"log"
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
	updatedPRBuffer, _, err := gh.Exec("search", "prs", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "--json", "number", "--jq", ".[].number")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Starting to collect the file changes from PR's")
	scanner := bufio.NewScanner(&updatedPRBuffer)
	for scanner.Scan() {
		getPRDetails(scanner.Text())
		//fmt.Printf("metric: %s", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func getPRDetails(prNumber string) {
	prDetails, _, err := gh.Exec("pr", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "diff", prNumber)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(prDetails.String())
	//fmt.Println(_prDetails)
}

package main

import (
	"encoding/json"
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
	updatedPR, _, err := gh.Exec("search", "prs", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "--json", "number")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(updatedPR.String())
	//var r Result
	//if err := json.Unmarshal(updatedPR.Bytes(), &r); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(r.Branch)
	//splitString := strings.Split(updatedPR.String(), ",")
	//for prNumber := range splitString {
	//	getPRDetails(prNumber)
	//}

}

//func getPRDetails(prNumber int) {
//	gh.Exec("pr", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "diff", string(prNumber))
//}

// Will add constants once on work PC
const (
	Summer string = "summer"
	Autumn        = "autumn"
	Winter        = "winter"
	Spring        = "spring"
)

func (r *Result) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &r.PRNum); err != nil {
		return err
	}

	if len(tmp) > 1 {
		if err := json.Unmarshal(tmp[1], &r.Branch); err != nil {
			return err
		}
		if err := json.Unmarshal(tmp[2], &r.FilesChanged); err != nil {
			return err
		}
		if err := json.Unmarshal(tmp[3], &r.Changes); err != nil {
			return err
		}
	}
	return nil
}

type Result struct {
	PRNum        string
	Branch       string
	FilesChanged string
	Changes      string
}

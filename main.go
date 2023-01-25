package main

import (
	"fmt"

	"github.com/cli/go-gh"
)

// This is a test comment
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
	/*	updatedPRBuffer, _, err := gh.Exec("search", "prs", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io", "--json", "number", "--jq", ".[].number")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(updatedPRBuffer)

		prListResult := &Results{}
		err = json.Unmarshal(updatedPRBuffer.Bytes(), prListResult)
		if err != nil {
			fmt.Println(err)
		}*/

	tempTest, _, err := gh.Exec("api", "graphql", "-F", "$org", "Personal-Development-Projects", "-F", "$repo", "OConnor-Development-Project.github.io", "-F", "query", PrListQuery)
	if err != nil {
		return
	}
	fmt.Println(tempTest)

	//err = json.Unmarshal(updatedPR.Bytes(), &prListResult)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println(updatedPR.Bytes())
	//	return
	//}

	//if err := json.Unmarshal(updatedPR.Bytes(), &prListResult); err != nil {
	//	log.Fatal(err)
	//}
	// gh api \
	//  -H "Accept: application/vnd.github+json" \
	//  /repos/OWNER/REPO/pulls/PULL_NUMBER
	// USE FOR Getting PR number with
	//updatedP, _, err := gh.Exec("api", "Accept", "application/vnd.github.diff+json", "--repo", "Personal-Development-Projects/OConnor-Development-Project.github.io/pulls/1")

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

//func (result *PRListResult) UnmarshalJSON(p []byte) error {
//	var tmp []json.RawMessage
//	if err := json.Unmarshal(p, &tmp); err != nil {
//		return err
//	}
//	if err := json.Unmarshal(tmp[0], &result.Number); err != nil {
//		return err
//	}
//
//	if len(tmp) > 1 {
//		if err := json.Unmarshal(tmp[1], &result.Branch); err != nil {
//			return err
//		}
//		if err := json.Unmarshal(tmp[2], &result.FilesChanged); err != nil {
//			return err
//		}
//		if err := json.Unmarshal(tmp[3], &result.Changes); err != nil {
//			return err
//		}
//	}
//	return nil
//}

type Results struct {
	_ []PullRequest
}

// Result base object for PRs

type PullRequest struct {
	number string
}

const (
	PrListQuery string = "'\n query allPullRequests($org: String!, $repo: String!, $endCursor: String) {\n  organization(login: $org) {\n    repository(name: $repo) {\n      pullRequests(first: 100, after: $endCursor, states: OPEN,) {\n nodes {\n author {\n login\n }\n number\n createdAt\n mergedAt\n mergedBy {\n login\n }\n approvers: reviews(states: APPROVED, first: 10) {\n nodes {\n author {\n ... on User {\n name\n login\n }\n }\n state\n }\n }\n title\n repository {\n owner {\n login\n }\n name\n }\n }\n pageInfo {\n hasNextPage\n endCursor\n }\n }\n }\n }\n}\n'"
	PrNumQuery  string = "summer"
)

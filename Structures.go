package main

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

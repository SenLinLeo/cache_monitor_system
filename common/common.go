package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// *IssuesSearchResult
func SearchIssues(issueURL string, result interface{}) (error) {

	resp, err := http.Get(issueURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("search query failed: Status:[%s]", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return  err
	}

	return  nil
}

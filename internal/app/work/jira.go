package work

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"os"
)

type CodeBase string

const (
	Mobile   CodeBase = "mobile"
	Frontend CodeBase = "frontend"
	Backend  CodeBase = "backend"
)

type JiraWork struct {
	Client *jira.Client
}

func NewJira() *JiraWork {
	jt := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}

	client, err := jira.NewClient(jt.Client(), os.Getenv("JIRA_URL"))
	if err != nil {
		fmt.Println(err)
	}

	return &JiraWork{
		Client: client,
	}
}

func (j *JiraWork) GetFirstIssueByCodeBase(base CodeBase) *jira.Issue {
	opt := &jira.SearchOptions{
		MaxResults: 1,
	}

	jql := fmt.Sprintf("project = \"CD\" AND status = Backlog AND \"codebase[dropdown]\" = %s ORDER BY created ASC", base)
	issues, _, err := j.Client.Issue.Search(jql, opt)
	if err != nil {
		fmt.Println(err)
	}

	// Return the first issue
	if len(issues) == 0 {
		return nil
	}

	return &issues[0]
}

package interfaces

import (
	"github.com/andygrunwald/go-jira"
	"github.com/manomartins/bitbird/internal/app/work"
)

type IssueService interface {
	GetFirstIssueByCodeBase(base work.CodeBase) *jira.Issue
}

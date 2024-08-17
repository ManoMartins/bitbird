package events

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/manomartins/bitbird/internal/interfaces"
	"github.com/manomartins/bitbird/internal/model"
	"github.com/manomartins/bitbird/internal/work"
	"regexp"
	"sync"
)

var wg sync.WaitGroup

type CheckCD struct {
	notifier               interfaces.Notifier
	issueService           interfaces.IssueService
	deploymentQueueStorage interfaces.DeploymentQueueInterface
}

func NewCheckCD(notifier interfaces.Notifier, issueService interfaces.IssueService, deploymentQueueStorage interfaces.DeploymentQueueInterface) *CheckCD {
	return &CheckCD{
		issueService:           issueService,
		notifier:               notifier,
		deploymentQueueStorage: deploymentQueueStorage,
	}
}

func (c *CheckCD) Execute() error {
	issues := make(chan *jira.Issue, 3)

	wg.Add(3)

	go c.getFirstIssueByCodeBase(work.Mobile, issues)
	go c.getFirstIssueByCodeBase(work.Backend, issues)
	go c.getFirstIssueByCodeBase(work.Frontend, issues)

	wg.Wait()
	close(issues)

	for issue := range issues {
		if issue == nil {
			continue
		}

		hash, err := c.extractHash(issue.Fields.Summary)
		if err != nil {
			return err
		}

		message := c.generateDeployNotification(
			issue.Key,
			issue.Fields.Assignee.DisplayName,
			hash,
		)

		messageID, err := c.notifier.SendNotification(message)
		if err != nil {
			return err
		}

		err = c.deploymentQueueStorage.Create(model.DeploymentQueueModel{
			CardKey:    issue.Key,
			MessageID:  messageID,
			CommitHash: hash,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CheckCD) getFirstIssueByCodeBase(base work.CodeBase, issueChan chan *jira.Issue) {
	defer wg.Done()
	issue := c.issueService.GetFirstIssueByCodeBase(base)

	if issue == nil || issue.Fields.Assignee == nil {
		issueChan <- nil
		return
	}

	key, err := c.deploymentQueueStorage.GetByCardKey(issue.Key)
	if err != nil {
		issueChan <- nil
		return
	}

	if key != nil {
		issueChan <- nil
		return
	}

	issueChan <- issue
}

func (c *CheckCD) extractHash(deploymentString string) (string, error) {
	// Define the regular expression to capture the commit hash
	re := regexp.MustCompile(`[a-fA-F0-9]{40}`)

	// Find the first occurrence that matches the pattern
	hash := re.FindString(deploymentString)

	if hash == "" {
		return "", fmt.Errorf("hash not found in the provided string")
	}

	return hash, nil
}

func (c *CheckCD) generateDeployNotification(cardKey, author, hash string) string {
	message := "🔔 **Deploy em Homologação**\n\n"
	message += fmt.Sprintf("**Chave do Card:** %s\n", cardKey)
	message += fmt.Sprintf("**Autor:** %s\n", author)
	message += fmt.Sprintf("**Hash do Commit:** %s\n\n", hash)
	message += fmt.Sprintf("📋 **Comando Git:**\n")
	message += fmt.Sprintf("```bash\n")
	message += fmt.Sprintf("git checkout -b %s origin/homolog\n", cardKey)
	message += fmt.Sprintf("git cherry-pick %s\n", hash)
	message += fmt.Sprintf("```\n")
	message += fmt.Sprintf("⚠️ **Ação:** Realizar o cherry-pick usando o comando acima e revisar o código.\n")

	return message
}

package storage

import (
	"encoding/json"
	"fmt"
	"github.com/manomartins/bitbird/internal/model"
	"io/ioutil"
	"log"
	"os"
)

type PullRequestMessages struct {
}

var pullRequestMessages []model.PullRequestMessageModel

func NewPullRequestMessages() *PullRequestMessages {
	return &PullRequestMessages{}
}

func (p *PullRequestMessages) GetPullRequestMessage(prID string) (*model.PullRequestMessageModel, error) {
	for _, message := range pullRequestMessages {
		if message.PrID == prID {
			return &message, nil
		}
	}

	return nil, fmt.Errorf("PR with ID %s not found", prID)
}

func (p *PullRequestMessages) FindAllPullRequestMessages() ([]model.PullRequestMessageModel, error) {
	file, err := os.Open("pull_request_messages.json")
	if err != nil {
		return nil, fmt.Errorf("No existing mappings file found, creating a new one.")
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, &pullRequestMessages)
	if err != nil {
		return nil, err
	}

	return pullRequestMessages, nil
}

func (p *PullRequestMessages) SavePullRequestMessage() {
	file, err := os.Create("pull_request_messages.json")
	if err != nil {
		log.Fatalf("Error creating JSON file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(pullRequestMessages)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
}

func (p *PullRequestMessages) UpdatePullRequestMessage(prID string, messageID string) error {
	for i, message := range pullRequestMessages {
		if message.PrID == prID {
			pullRequestMessages[i].MessageID = messageID
			p.SavePullRequestMessage()
			return nil
		}
	}

	pullRequestMessages = append(pullRequestMessages, model.PullRequestMessageModel{PrID: prID, MessageID: messageID})
	p.SavePullRequestMessage()

	return nil
}

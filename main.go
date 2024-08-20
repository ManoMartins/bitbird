package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/manomartins/bitbird/configs"
	chat_notify "github.com/manomartins/bitbird/internal/chat_notify"
	events "github.com/manomartins/bitbird/internal/events"
	"github.com/manomartins/bitbird/internal/interfaces"
	storage "github.com/manomartins/bitbird/internal/storage"
	"github.com/manomartins/bitbird/internal/work"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
)

var notifier interfaces.Notifier

type eventFunc func(event events.PullRequestEvent) error

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	var event events.PullRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Cannot unmarshal JSON", http.StatusBadRequest)
		return
	}

	eventKey := r.Header.Get("X-Event-Key")

	messagesStorage := storage.NewPullRequestMessagesMongo()
	eventPullRequestCreated := events.NewPullRequestCreated(notifier, messagesStorage)
	eventPullRequestApproved := events.NewPullRequestApproved(notifier, messagesStorage)
	eventPullRequestChangesRequest := events.NewPullRequestChangesRequest(notifier, messagesStorage)
	eventPullRequestCommentCreated := events.NewPullRequestCommentCreated(notifier, messagesStorage)
	eventPullRequestRemovedAction := events.NewPullRequestRemovedAction(notifier, messagesStorage)

	eventHandlers := map[string]eventFunc{
		"pullrequest:created":                 eventPullRequestCreated.Execute,
		"pullrequest:approved":                eventPullRequestApproved.Execute,
		"pullrequest:changes_request_created": eventPullRequestChangesRequest.Execute,
		"pullrequest:comment_created":         eventPullRequestCommentCreated.Execute,
		"pullrequest:unapproved":              eventPullRequestRemovedAction.Execute,
		"pullrequest:changes_request_removed": eventPullRequestRemovedAction.Execute,
	}

	if handler, exists := eventHandlers[eventKey]; exists {
		err = handler(event)

		fmt.Println("Error handling event", err)

		if err != nil {
			http.Error(w, "Error handling event", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Fatal(".env file missing")
	}
}

func main() {
	storage.ConnectMongoDB()
	defer storage.CloseMongoDB()

	configs.ConnectDiscord()
	defer configs.DisconnectDiscord()

	notifier = chat_notify.NewDiscordNotifier()

	c := cron.New()
	c.AddFunc("*/60 * * * *", func() {
		fmt.Println("Check CD")
		jiraWork := work.NewJira()
		deploymentQueue := storage.NewDeploymentQueueMongo()
		checkCD := events.NewCheckCD(notifier, jiraWork, deploymentQueue)

		err := checkCD.Execute()
		if err != nil {
			return
		}
	})
	c.Start()

	http.HandleFunc("/", handleWebhook)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	log.Println("Listening on :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

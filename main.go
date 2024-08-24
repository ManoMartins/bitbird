package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/joho/godotenv"
	"github.com/manomartins/bitbird/configs"
	chat_notify "github.com/manomartins/bitbird/internal/chat_notify"
	events "github.com/manomartins/bitbird/internal/events"
	"github.com/manomartins/bitbird/internal/interfaces"
	storage "github.com/manomartins/bitbird/internal/storage"
	"github.com/manomartins/bitbird/internal/work"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var notifier interfaces.Notifier

type eventFunc func(ctx context.Context, event events.PullRequestEvent) error

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	var event events.PullRequestEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Cannot unmarshal JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	eventKey := r.Header.Get("X-Event-Key")

	span := trace.SpanFromContext(ctx)
	span.AddEvent("pull_request_created.event")
	span.SetAttributes(
		attribute.String("pull_request.event_key", eventKey),
	)

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
		err = handler(ctx, event)

		if err != nil {
			log.Println("Error handling event", err)
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
	ctx := context.Background()

	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	storage.ConnectMongoDB()
	defer storage.CloseMongoDB()

	configs.ConnectDiscord()
	defer configs.DisconnectDiscord()

	notifier = chat_notify.NewDiscordNotifier()

	c := cron.New()
	c.AddFunc("*/30 10-23 * * *", func() {
		log.Println("Checking deployment queue...")
		jiraWork := work.NewJira()
		deploymentQueue := storage.NewDeploymentQueueMongo()
		checkCD := events.NewCheckCD(notifier, jiraWork, deploymentQueue)

		err := checkCD.Execute(ctx)
		if err != nil {
			return
		}
	})
	c.Start()

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(handleWebhook), "webhooks")

	http.Handle("/", otelHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Listening on :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

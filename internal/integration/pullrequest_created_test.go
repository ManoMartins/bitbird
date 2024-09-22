package integration

import (
	"context"
	"fmt"
	"github.com/manomartins/bitbird/internal/app/events"
	"github.com/manomartins/bitbird/internal/app/storage"
	mocks "github.com/manomartins/bitbird/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strconv"
	"testing"
)

// setupMongoDB inicializa um container do MongoDB e configura a conexão.
// Retorna o contexto, uma função de cleanup para liberar os recursos e um erro, caso ocorra.
func setupMongoDB(t *testing.T) (context.Context, func(), error) {
	ctx := context.TODO()

	// Inicializa o container do MongoDB
	mongodbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
		},
	})

	if err != nil {
		t.Fatalf("failed to start MongoDB container: %s", err)
	}

	// Pega o host e a porta mapeada do container
	host, err := mongodbContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := mongodbContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}

	// Configura as variáveis de ambiente necessárias para a conexão com o MongoDB
	os.Setenv("MONGO_USER", "root")
	os.Setenv("MONGO_PASS", "root")
	os.Setenv("MONGO_HOST", fmt.Sprintf("%s:%s", host, port.Port()))

	// Conecta ao MongoDB
	storage.ConnectMongoDB()

	// Função de cleanup que deve ser chamada após o teste para limpar os recursos
	cleanup := func() {
		storage.CloseMongoDB() // Fecha a conexão com o banco
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate MongoDB container: %s", err)
		}
	}

	return ctx, cleanup, nil
}

func TestPullRequestCreated_Execute(t *testing.T) {
	os.Setenv("DISCORD_CHANNEL_ID_FOR_PR", "channel-123")

	// Configura o MongoDB para o teste
	ctx, cleanup, err := setupMongoDB(t)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup() // Garante que o cleanup seja executado ao final do teste

	// Mock de notificações
	notifierMock := new(mocks.MockNotifier)
	messagesStorage := storage.NewPullRequestMessagesMongo()

	// Cria o evento de Pull Request
	pullRequestCreated := events.NewPullRequestCreated(notifierMock, messagesStorage)

	notifierMock.
		On("GetUserAvatarURL", mock.Anything, "667184274428002345").
		Return("https://cdn.discordapp.com/avatars/667184274428002345/avatar.png", nil)
	notifierMock.
		On("SendNotificationEmbed", mock.Anything, mock.Anything, mock.Anything).
		Return("message-123", nil)

	event := events.PullRequestEvent{
		Actor: events.Actor{
			DisplayName: "manoel_martins",
		},
		PullRequest: events.PullRequest{
			ID:    1,
			Title: "Adicionar nova funcionalidade",
			State: "OPEN",
			Destination: events.Destination{
				Branch: events.Branch{Name: "main"},
			},
			Reviewers: []events.Reviewer{
				{DisplayName: "alexandre_valim"},
			},
			Links: events.PullRequestLinks{
				HTML: events.HTML{Href: "https://github.com/repo/pull/1"},
			},
		},
		Repository: events.Repository{Name: "bitbird"},
	}

	// Executa o evento
	err = pullRequestCreated.Execute(ctx, event)
	assert.NoError(t, err)

	// Verifica se as expectativas do mock foram atendidas
	notifierMock.AssertExpectations(t)

	// Verifica se a mensagem foi armazenada corretamente
	message, err := messagesStorage.GetById(strconv.Itoa(event.PullRequest.ID))
	assert.NoError(t, err)
	assert.Equal(t, message.PrID, strconv.Itoa(event.PullRequest.ID))
	assert.Equal(t, message.ChannelID, "channel-123")
	assert.Equal(t, message.MessageID, "message-123")
}

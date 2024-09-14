package integration

import (
	"context"
	"github.com/manomartins/bitbird/internal/app/events"
	mocks "github.com/manomartins/bitbird/internal/mocks"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

func TestPullRequestCreated_Execute(t *testing.T) {
	// Configuração do ambiente de teste com Testcontainers (Redis como exemplo)
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:5.0.3-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})

	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// Clean up the container
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	// Mocks
	notifierMock := new(mocks.MockNotifier)
	messagesStorageMock := new(mocks.MockPullRequestMessagesInterface)

	// Criar a instância de PullRequestCreated com os mocks
	pullRequestCreated := events.NewPullRequestCreated(notifierMock, messagesStorageMock)

	// Configurar o mock do notifier para retornar um avatar URL válido
	notifierMock.On("GetUserAvatarURL", mock.Anything, "667184274428002345").Return("https://cdn.discordapp.com/avatars/667184274428002345/avatar.png", nil)
	notifierMock.On("SendNotificationEmbed", mock.Anything, mock.Anything, mock.Anything).Return("123456789", nil)

	// Configurar o mock do messagesStorage para salvar a mensagem
	messagesStorageMock.On("Create", mock.Anything, "1", mock.Anything, "123456789").Return(nil)

	// Dados do evento para o teste
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
			Links: events.Links{
				HTML: events.Link{Href: "https://github.com/repo/pull/1"},
			},
		},
		Repository: events.Repository{Name: "bitbird"},
	}

	// Executar o método que está sendo testado
	err = pullRequestCreated.Execute(ctx, event)

	// Verificar se não houve erro
	assert.NoError(t, err)

	// Verificar se as funções do mock foram chamadas corretamente
	notifierMock.AssertExpectations(t)
	messagesStorageMock.AssertExpectations(t)
}

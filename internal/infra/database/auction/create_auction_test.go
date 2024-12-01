package auction_test

import (
	"github.com/emebit/goexperts-lab-auction/internal/entity/auction_entity"
	"github.com/emebit/goexperts-lab-auction/internal/infra/database/auction"
	"context"
	"testing"
	"time"
)

func TestCreateAuction(t *testing.T) {
	mockRepo := auction.NewAuctionRepositoryMock()

	auctionEntity := &auction_entity.Auction{
		Id:          "Leilao do teste",
		ProductName: "Produto do Teste",
		Category:    "Categoria do Teste",
		Description: "Descrição do Teste",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	// Simulando a criação do leilão
	err := mockRepo.CreateAuction(context.Background(), auctionEntity)
	if err != nil {
		t.Fatalf("Failed to create auction: %v", err)
	}

	// Verificar se o leilão foi criado corretamente
	createdAuction, err := mockRepo.FindAuctionById(context.Background(), "Leilao do teste")
	if err != nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if createdAuction.ProductName != "Produto do Teste" {
		t.Errorf("Expected product name to be 'Produto do Teste', got %s", createdAuction.ProductName)
	}
}

func TestMonitorAuctions_AlreadyCompleted(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar leilão já encerrado
	now := time.Now()
	completedAuction := &auction_entity.Auction{
		Id:          "Leilao_1",
		ProductName: "Produto_1",
		Status:      auction_entity.Completed,
		Timestamp:   now.Add(-time.Hour), // Expirado há mais de uma hora
	}

	mockRepo.SaveAuction(completedAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAuctions(context.Background())

	// Verificar se o leilão já encerrado permanece encerrado
	completedAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "Leilao_1")
	if err != nil {
		t.Fatalf("Error finding completed auction: %v", err)
	}
	if completedAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected completed auction status to remain 'Completed', got %v", completedAuctionResult.Status)
	}
}

func TestMonitorAuctions_NotExpired(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar leilão que não deve ser encerrado
	now := time.Now()
	notExpiredAuction := &auction_entity.Auction{
		Id:          "Leilao_2",
		ProductName: "Produto_2",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Second * 10), // Não Expirado, menos de 20 segundos
	}

	mockRepo.SaveAuction(notExpiredAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAuctions(context.Background())

	// Verificar se o leilão não expirado ainda está ativo
	notExpiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "Leilao_2")
	if err != nil {
		t.Fatalf("Error finding not expired auction: %v", err)
	}
	if notExpiredAuctionResult.Status != auction_entity.Active {
		t.Errorf("Expected not expired auction status to be 'Active', got %v", notExpiredAuctionResult.Status)
	}
}

func TestMonitorAuctions_MixedStatus(t *testing.T) {
	// Configuração inicial
	mockRepo := auction.NewAuctionRepositoryMock()

	// Criar alguns leilões simulados com status diferentes
	now := time.Now()
	expiredAuction := &auction_entity.Auction{
		Id:          "Leilao_1",
		ProductName: "Produto_1",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Minute * 40),
	}
	notExpiredAuction := &auction_entity.Auction{
		Id:          "Leilao_2",
		ProductName: "Produto_2",
		Status:      auction_entity.Active,
		Timestamp:   now.Add(-time.Second * 10),
	}
	completedAuction := &auction_entity.Auction{
		Id:          "Leilao_3",
		ProductName: "Produto_3",
		Status:      auction_entity.Completed,
		Timestamp:   now.Add(-time.Hour),
	}

	mockRepo.SaveAuction(expiredAuction)
	mockRepo.SaveAuction(notExpiredAuction)
	mockRepo.SaveAuction(completedAuction)

	// Monitorar e fechar leilões expirados
	mockRepo.MonitorAuctions(context.Background())

	// Verificar se o leilão expirado foi fechado
	expiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "Leilao_1")
	if err != nil {
		t.Fatalf("Error finding expired auction: %v", err)
	}
	if expiredAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected expired auction status to be 'Completed', got %v", expiredAuctionResult.Status)
	}

	// Verificar se o leilão não expirado ainda está ativo
	notExpiredAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "Leilao_2")
	if err != nil {
		t.Fatalf("Error finding not expired auction: %v", err)
	}
	if notExpiredAuctionResult.Status != auction_entity.Active {
		t.Errorf("Expected not expired auction status to be 'Active', got %v", notExpiredAuctionResult.Status)
	}

	// Verificar se o leilão já encerrado permanece encerrado
	completedAuctionResult, err := mockRepo.FindAuctionById(context.Background(), "Leilao_3")
	if err != nil {
		t.Fatalf("Error finding completed auction: %v", err)
	}
	if completedAuctionResult.Status != auction_entity.Completed {
		t.Errorf("Expected completed auction status to remain 'Completed', got %v", completedAuctionResult.Status)
	}
}

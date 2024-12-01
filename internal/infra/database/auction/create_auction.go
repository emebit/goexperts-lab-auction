package auction

import (
	"github.com/emebit/goexperts-lab-auction/configuration/logger"
	"github.com/emebit/goexperts-lab-auction/internal/entity/auction_entity"
	"github.com/emebit/goexperts-lab-auction/internal/internal_error"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	cancelFunc context.CancelFunc //Adiconado
}

// Estrutura para armazenar dados do filtro leilão expirado
type expiredAuctionsFilter struct {
	Status    auction_entity.AuctionStatus `bson:"status"`
	Timestamp primitive.M                  `bson:"timestamp"`
}

// Estrutura de Status Completo
type updateCompletedStatus struct {
	SetStatus auction_entity.AuctionStatus `bson:"$set.status"`
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	ctx, cancel := context.WithCancel(context.Background()) //Adicionado
	ar := &AuctionRepository{
		Collection: database.Collection("auctions"),
		cancelFunc: cancel,
	}
	//Cria thread chamando funcção para monitorar os leilões
	go ar.MonitorAuctions(ctx) //Adicionado
	return ar                  //Alterado
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

// Função para pegar o valor do parametro de duração
func getAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	interval, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Second * 30 //Se não conseguir parsear auctionInterval, arbitra 30 segundos
	}

	return interval
}

// Função que encerra o contexto
func (ar *AuctionRepository) Cancel() {
	ar.cancelFunc()
}

// Função que cria filtro de leilão expirado
func (ar *AuctionRepository) createExpiredAuctionsFilter(auctionInterval time.Duration, now int64) expiredAuctionsFilter {
	return expiredAuctionsFilter{
		Status:    auction_entity.Active,
		Timestamp: primitive.M{"$lt": now - int64(auctionInterval.Seconds())},
	}
}

// Função que cria o Status de Completo
func (ar *AuctionRepository) createUpdateCompletedStatus() updateCompletedStatus {
	return updateCompletedStatus{
		SetStatus: auction_entity.Completed,
	}
}

// Função que Atualiza os leilões expirados
func (ar *AuctionRepository) updateExpiredAuctions(ctx context.Context, filter expiredAuctionsFilter, update updateCompletedStatus) {
	_, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Erro ao tentar atualizar leilão expirado", err)
	}
}

// Função para verificar e fechar leilão expirado
func (ar *AuctionRepository) checkExpiredAuctions(ctx context.Context) {
	auctionInterval := getAuctionInterval()
	now := time.Now().Unix()

	filter := ar.createExpiredAuctionsFilter(auctionInterval, now)
	update := ar.createUpdateCompletedStatus()

	ar.updateExpiredAuctions(ctx, filter, update)
}

// Função para monitorar leilão expirado
func (ar *AuctionRepository) MonitorAuctions(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Info("Verificando leilão expirado")
			ar.checkExpiredAuctions(ctx)
		case <-ctx.Done():
			logger.Info("Parando goroutine MonitorAuctions")
			return
		}
	}
}

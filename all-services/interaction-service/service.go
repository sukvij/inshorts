package interactionservice

import (
	"context"

	"github.com/go-redis/redis"
	newsservice "github.com/sukvij/inshorts/all-services/news-service"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type InteractionService struct {
	DB           *gorm.DB
	Interactions *[]UserInteraction
	Redis        *redis.Client
}

func _NewService(db *gorm.DB, interatcions *[]UserInteraction, redis *redis.Client) *InteractionService {
	return &InteractionService{DB: db, Interactions: interatcions, Redis: redis}
}

func (service *InteractionService) CreateUserInteraction() error {
	repo := _NewRepository(service.DB, service.Interactions)
	return repo.CreateUserInteraction()
}

func (service *InteractionService) TrendingNewsArticles(lat, lon float64, limit int, cacheKey string, radiusMeter int) (*[]newsservice.NewsArticle, error) {

	_, span := otel.Tracer("function controller").Start(context.Background(), "service redis data fetch")
	defer span.End()

	repo := _NewRepository(service.DB, service.Interactions)
	res, errs := repo.TrendingNewsArticles(lat, lon, limit, radiusMeter)
	return res, errs
}

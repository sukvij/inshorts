package newsservice

import "gorm.io/gorm"

type NewsService struct {
	DB          *gorm.DB
	NewsArticle *[]NewsArticle
}

func _NewService(db *gorm.DB, newsArticles *[]NewsArticle) *NewsService {
	return &NewsService{DB: db, NewsArticle: newsArticles}
}

func (service *NewsService) CreateNewsArticle() error {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.CreateNewsArticle()
}

func (service *NewsService) GetNewsArticlesByCategory(category string) (*[]NewsArticle, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesByCategory(category)
}

func (service *NewsService) GetNewsArticlesByScore(score float64) (*[]NewsArticle, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesByScore(score)
}

func (service *NewsService) GetNewsArticlesBySource(source string) (*[]NewsArticle, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesBySource(source)
}

func (service *NewsService) GetNearByNewsArticle(lat, lon, radius float64) (*[]NewsArticle, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNearByNewsArticle(lat, lon, radius)
}

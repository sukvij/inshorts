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

func (service *NewsService) GetNewsArticles() (*[]NewsArticle, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticles()
}

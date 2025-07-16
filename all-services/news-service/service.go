package newsservice

import (
	"fmt"
	"strings"

	llmservice "github.com/sukvij/inshorts/all-services/llm-service"
	"gorm.io/gorm"
)

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

func (service *NewsService) GetNewsArticlesByCategory(category string) (*[]NewsArticle, int64, string, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesByCategory(category)
}

func (service *NewsService) GetNewsArticlesByScore(score float64) (*[]NewsArticle, int64, string, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesByScore(score)
}

func (service *NewsService) GetNewsArticlesBySource(source string) (*[]NewsArticle, int64, string, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticlesBySource(source)
}

func (service *NewsService) GetNearByNewsArticle(lat, lon, radius float64) (*[]NewsArticle, int64, string, error) {
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNearByNewsArticle(lat, lon, radius)
}

func (service *NewsService) GetNewsArticleBySearch(llmOutput *llmservice.ParsedLLMOutput) (*[]NewsArticle, int64, string, error) {
	// funny business
	fmt.Println("llmOutput.Entities ", llmOutput.Entities)
	searchTerms := []string{}
	for _, term := range strings.Split(llmOutput.Entities, ",") {
		cleanedTerm := strings.TrimSpace(term)
		if cleanedTerm != "" {
			searchTerms = append(searchTerms, strings.ToLower(cleanedTerm)) // Convert to lowercase for case-insensitive search
		}
	}

	if len(searchTerms) == 0 {
		return nil, 0, "", fmt.Errorf("empty result") // No search terms, return empty slice
	}

	// (LOWER(title) LIKE '%term1%' OR LOWER(description) LIKE '%term1%') OR
	// var conditions []string
	var args []interface{}

	// for _, term := range searchTerms {
	// 	conditions = append(conditions, "(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)")
	// 	args = append(args, "%"+term+"%", "%"+term+"%")
	// }
	// for _, term := range searchTerms {
	// 	conditions = append(conditions, term)
	// 	// args = append(args, "%"+term+"%", "%"+term+"%")
	// }
	whereClause := strings.Join(searchTerms, " ")
	fmt.Println("where clause if ", whereClause)
	repo := _NewRepository(service.DB, service.NewsArticle)
	return repo.GetNewsArticleBySearch(whereClause, args)
}

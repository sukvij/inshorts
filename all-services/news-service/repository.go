package newsservice

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type NewsRepository struct {
	DB          *gorm.DB
	NewsArticle *[]NewsArticle
}

func _NewRepository(db *gorm.DB, news *[]NewsArticle) *NewsRepository {
	return &NewsRepository{DB: db, NewsArticle: news}
}

func (repo *NewsRepository) CreateNewsArticle() error {

	for _, article := range *repo.NewsArticle {
		err := repo.DB.Create(&article).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *NewsRepository) GetNewsArticlesByCategory(category string) (*[]NewsArticle, error) {
	marshaledCategory, err1 := json.Marshal(category)
	if err1 != nil {
		return nil, err1
	}
	var result []NewsArticle

	// Use GORM's Where clause with MySQL's JSON_CONTAINS function.
	// JSON_CONTAINS(json_doc, val) returns 1 if val is found in json_doc.
	err := repo.DB.Where("JSON_CONTAINS(category, ?)", string(marshaledCategory)).Order("publication_date desc").Find(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *NewsRepository) GetNewsArticlesByScore(score float64) (*[]NewsArticle, error) {
	var result []NewsArticle
	err := repo.DB.Where("relevance_score > ?", score).Order("relevance_score desc").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *NewsRepository) GetNewsArticlesBySource(source string) (*[]NewsArticle, error) {
	var result []NewsArticle
	err := repo.DB.Where("source_name = ?", source).Order("publication_date desc").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *NewsRepository) GetNearByNewsArticle(lat, lon, radius float64) (*[]NewsArticle, error) {
	var result []NewsArticle

	query := fmt.Sprintf(`
								SELECT
								*,
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS(%v)) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS(%v)) +
										SIN(RADIANS(%v)) * SIN(RADIANS(latitude))
									)
								) AS distance 
							FROM
								news_articles
							WHERE
								(
									6371 * -- Earth's radius in kilometers
									ACOS(
										COS(RADIANS(%v)) * COS(RADIANS(latitude)) *
										COS(RADIANS(longitude) - RADIANS(%v)) +
										SIN(RADIANS(%v)) * SIN(RADIANS(latitude))
									)
								) <= %v 
							ORDER BY
								distance 
							LIMIT 2;
								
								`, lat, lon, lat, lat, lon, lat, radius)

	err := repo.DB.Raw(query).Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch points of interest within radius: %w", err)
	}

	return &result, nil
}

func (repo *NewsRepository) GetNewsArticleBySearch(whereClause string, arg []interface{}) (*[]NewsArticle, error) {
	var result []NewsArticle
	err := repo.DB.Where(whereClause, arg...).Find(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles from DB: %v", err)
	}

	return &result, nil
}
